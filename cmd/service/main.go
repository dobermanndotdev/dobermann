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

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/flowck/dobermann/backend/internal/adapters/endpoint_checkers"
	"github.com/flowck/dobermann/backend/internal/adapters/events"
	"github.com/flowck/dobermann/backend/internal/adapters/psql"
	"github.com/flowck/dobermann/backend/internal/adapters/resend"
	"github.com/flowck/dobermann/backend/internal/adapters/transaction"
	"github.com/flowck/dobermann/backend/internal/app"
	"github.com/flowck/dobermann/backend/internal/app/command"
	"github.com/flowck/dobermann/backend/internal/app/query"
	"github.com/flowck/dobermann/backend/internal/common/auth"
	"github.com/flowck/dobermann/backend/internal/common/logs"
	"github.com/flowck/dobermann/backend/internal/common/messaging"
	"github.com/flowck/dobermann/backend/internal/common/observability"
	"github.com/flowck/dobermann/backend/internal/common/postgres"
	"github.com/flowck/dobermann/backend/internal/common/watermill_logger"
	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
	amqpport "github.com/flowck/dobermann/backend/internal/ports/amqp"
	httpport "github.com/flowck/dobermann/backend/internal/ports/http"
)

var Version = "development"

type Config struct {
	AmqpUrl                       string `envconfig:"AMQP_URL"`
	Port                          int    `envconfig:"HTTP_PORT"`
	JwtSecret                     string `envconfig:"JWT_SECRET"`
	DebugMode                     string `envconfig:"DEBUG_MODE"`
	DatabaseURL                   string `envconfig:"DATABASE_URL"`
	ResendApiKey                  string `envconfig:"RESEND_API_KEY"`
	HostnameForNotifications      string `envconfig:"HOSTNAME_NOTIFICATION"`
	SentFromEmailAddress          string `envconfig:"SENT_FROM_EMAIL_ADDRESS"`
	IsProductionMode              bool   `envconfig:"PRODUCTION_MODE"`
	Region                        string `envconfig:"WORKER_REGION"`
	EndpointCheckTimeoutInSeconds int    `envconfig:"ENDPOINT_CHECK_TIMEOUT_IN_SECONDS" required:"true"`
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
	watermillLogger := watermill_logger.NewWatermillLogger(logger)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	db, err := postgres.Connect(config.DatabaseURL)
	if err != nil {
		logger.Fatal(err)
	}

	logger.Info("Connected successfully to the database")

	err = postgres.ApplyMigrations(db, "misc/sql/migrations")
	if err != nil {
		logger.Fatal(err)
	}

	tokenSigner, err := auth.NewTokenSigner(config.JwtSecret, time.Hour*24*7)
	if err != nil {
		logger.Fatal(err)
	}

	publisher, err := messaging.NewAmqpPublisher(config.AmqpUrl, logger)
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

		// Sets the correlation id to the messages' context
		messaging.CorrelationIdMiddleware,

		// Send failed events to a specific queue
		poisonQueueMiddleware,

		messaging.ErrorLoggerMiddleware(logger),

		// Retry failed events
		retryMiddleware.Middleware,
	)

	logger.Info("Connected successfully to RabbitMQ")

	httpChecker, err := endpoint_checkers.NewHttpChecker(config.Region, config.EndpointCheckTimeoutInSeconds)
	if err != nil {
		logger.Fatal(err)
	}

	userRepository := psql.NewUserRepository(db)
	eventPublisher := events.NewPublisher(publisher)
	monitorRepository := psql.NewMonitorRepository(db)
	txProvider := transaction.NewPsqlProvider(db, publisher, logger)

	//POST_IDEA: Mocking services and its dynamic initialisation
	var resendService resend.Service
	if config.IsProductionMode {
		//resendService = resend.NewService(config.ResendApiKey, config.SentFromEmailAddress, config.HostnameForNotifications)
		//logger.Info("Resend service has been initialised")
		resendService = resend.NewServiceMock(logger)
		logger.Warn("Resend service mock has been initialised")
	} else {
		resendService = resend.NewServiceMock(logger)
		logger.Warn("Resend service mock has been initialised")
	}

	application := &app.App{
		Commands: app.Commands{
			// Auth
			CreateAccount: observability.NewCommandDecorator[command.CreateAccount](command.NewCreateAccountHandler(txProvider), logger),
			LogIn:         observability.NewCommandWithResultDecorator[command.LogIn, string](command.NewLoginHandler(userRepository, tokenSigner), logger),

			// Monitor
			DeleteMonitor:                      observability.NewCommandDecorator[command.DeleteMonitor](command.NewDeleteMonitorHandler(txProvider), logger),
			CreateIncident:                     observability.NewCommandDecorator[command.CreateIncident](command.NewCreateIncidentHandler(txProvider), logger),
			EditMonitor:                        observability.NewCommandDecorator[command.EditMonitor](command.NewEditMonitorHandler(monitorRepository), logger),
			ResolveIncident:                    observability.NewCommandDecorator[command.ResolveIncident](command.NewResolveIncidentHandler(txProvider), logger),
			ToggleMonitorPause:                 observability.NewCommandDecorator[command.ToggleMonitorPause](command.NewToggleMonitorPauseHandler(txProvider), logger),
			CreateMonitor:                      observability.NewCommandDecorator[command.CreateMonitor](command.NewCreateMonitorHandler(monitorRepository, eventPublisher), logger),
			CheckEndpoint:                      observability.NewCommandDecorator[command.CheckEndpoint](command.NewCheckEndpointHandler(httpChecker, monitorRepository, eventPublisher), logger),
			NotifyMonitorSubscribersOnIncident: observability.NewCommandDecorator[command.NotifyMonitorSubscribersOnIncident](command.NewNotifyMonitorSubscribersOnIncidentHandler(txProvider, resendService), logger),
			NotifyOnIncidentResolved:           observability.NewCommandDecorator[command.NotifyOnIncidentResolved](command.NewNotifyOnIncidentResolvedHandler(monitorRepository, userRepository, resendService), logger),
		},
		Queries: app.Queries{
			MonitorByID:              observability.NewQueryDecorator[query.MonitorByID, *monitor.Monitor](query.NewMonitorByIdHandler(monitorRepository), logger),
			AllMonitors:              observability.NewQueryDecorator[query.AllMonitors, query.PaginatedResult[*monitor.Monitor]](query.NewAllMonitorsHandler(monitorRepository), logger),
			MonitorResponseTimeStats: observability.NewQueryDecorator[query.MonitorResponseTimeStats, []query.ResponseTimeStat](query.NewMonitorResponseTimeStatsHandler(monitorRepository), logger),

			// IAM
			UserByID: observability.NewQueryDecorator[query.UserByID, *account.User](query.NewUserByIdHandler(userRepository), logger),
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

	logger.WithFields(logs.Fields{
		"version": Version,
		"region":  config.Region,
	}).Info("The service is running")

	<-done
	terminationCtx, terminationCancel := context.WithTimeout(context.Background(), time.Second*2)
	defer terminationCancel()

	err = httpPort.Stop(terminationCtx)
	if err != nil {
		logger.Fatalf("unable to gracefully shutdown the http port: %v", err)
	}
}
