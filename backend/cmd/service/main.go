package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/flowck/dobermann/backend/internal/adapters/events"
	"github.com/flowck/dobermann/backend/internal/adapters/monitors"
	"github.com/flowck/dobermann/backend/internal/adapters/transaction"
	"github.com/flowck/dobermann/backend/internal/adapters/users"
	"github.com/flowck/dobermann/backend/internal/app"
	"github.com/flowck/dobermann/backend/internal/app/command"
	"github.com/flowck/dobermann/backend/internal/common/auth"
	"github.com/flowck/dobermann/backend/internal/common/logs"
	"github.com/flowck/dobermann/backend/internal/common/observability"
	"github.com/flowck/dobermann/backend/internal/common/psql"
	amqpport "github.com/flowck/dobermann/backend/internal/ports/amqp"
	httpport "github.com/flowck/dobermann/backend/internal/ports/http"
)

type Config struct {
	AmqpUrl     string `envconfig:"AMQP_URL"`
	Port        int    `envconfig:"HTTP_PORT"`
	JwtSecret   string `envconfig:"JWT_SECRET"`
	DebugMode   string `envconfig:"DEBUG_MODE"`
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func (c Config) IsDebugMode() bool {
	return strings.ToLower(c.DebugMode) == "enabled"
}

func main() {
	config := &Config{}
	err := envconfig.Process("", config)
	if err != nil {
		panic(err)
	}

	logger := logs.New(config.IsDebugMode())
	watermillLogger := watermill.NewStdLogger(false, false)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	db, err := psql.Connect(config.DatabaseURL)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connected successfully to the database")

	err = psql.ApplyMigrations(db, "misc/sql/migrations")
	if err != nil {
		logger.Fatal(err)
	}

	tokenSigner, err := auth.NewTokenSigner(config.JwtSecret, time.Hour*24*7)
	if err != nil {
		logger.Fatal(err)
	}

	publisher, err := amqp.NewPublisher(amqp.NewDurableQueueConfig(config.AmqpUrl), watermillLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() { _ = publisher.Close() }()

	subscriber, err := amqp.NewSubscriber(amqp.NewDurableQueueConfig(config.AmqpUrl), watermillLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() { _ = subscriber.Close() }()

	router, err := message.NewRouter(message.RouterConfig{}, watermillLogger)
	if err != nil {
		logger.Fatal(err)
	}
	defer func() { _ = router.Close() }()

	// TODO: https://github.com/ThreeDotsLabs/watermill/issues/173
	poisonQueueMiddleware, err := middleware.PoisonQueue(publisher, "poison_queue")
	if err != nil {
		logger.Fatal(err)
	}

	retryMiddleware := middleware.Retry{
		MaxRetries:      3,
		Logger:          watermillLogger,
		InitialInterval: time.Millisecond * 3,
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		// Handle panics
		middleware.Recoverer,

		// Get the correlation id
		middleware.CorrelationID,

		// Send failed events to a specific queue
		poisonQueueMiddleware,

		// Retry failed events
		retryMiddleware.Middleware,
	)

	logger.Info("Connected successfully to RabbitMQ")

	httpChecker := monitors.NewHttpChecker()
	userRepository := users.NewPsqlRepository(db)
	eventPublisher := events.NewPublisher(publisher)
	monitorRepository := monitors.NewPsqlRepository(db)
	incidentRepository := monitors.NewIncidentRepository(db)
	txProvider := transaction.NewPsqlProvider(db, publisher)

	application := &app.App{
		Commands: app.Commands{
			// Auth
			CreateAccount: observability.NewCommandDecorator[command.CreateAccount](command.NewCreateAccountHandler(txProvider), logger),
			LogIn:         observability.NewCommandWithResultDecorator[command.LogIn, string](command.NewLoginHandler(userRepository, tokenSigner), logger),

			// Monitor
			CreateIncident: observability.NewCommandDecorator[command.CreateIncident](command.NewCreateIncidentHandler(incidentRepository), logger),
			CreateMonitor:  observability.NewCommandDecorator[command.CreateMonitor](command.NewCreateMonitorHandler(monitorRepository, eventPublisher), logger),
			CheckEndpoint:  observability.NewCommandDecorator[command.CheckEndpoint](command.NewCheckEndpointHandler(httpChecker, monitorRepository, eventPublisher), logger),
		},
	}

	eventHandlers := amqpport.NewHandlers(application)
	for _, handler := range eventHandlers {
		router.AddNoPublisherHandler(handler.HandlerName(), handler.EventName(), subscriber, handler.Handle)
	}

	httpPort, err := httpport.NewPort(httpport.Config{
		Ctx:         ctx,
		Logger:      logger,
		Port:        config.Port,
		Application: application,
		JwtVerifier: tokenSigner,
	})
	if err != nil {
		logger.Fatal(err)
	}

	go func() {
		logger.Infof("The http port is running successfully at the port %d", config.Port)
		err = httpPort.Start()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Errorf("the http port stopped with the following error: %v", err)
		}
	}()

	go func() {
		err = router.Run(ctx)
		if err != nil {
			logger.Errorf("watermill router stopped with the following error: %v", err)
		}
	}()

	<-done
	terminationCtx, terminationCancel := context.WithTimeout(context.Background(), time.Second*2)
	defer terminationCancel()

	err = httpPort.Stop(terminationCtx)
	if err != nil {
		logger.Fatalf("unable to gracefully shutdown the http port: %v", err)
	}
}
