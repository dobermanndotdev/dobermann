package main

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/kelseyhightower/envconfig"
	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"
	resendsdk "github.com/resendlabs/resend-go"
	"github.com/sirupsen/logrus"

	"github.com/flowck/doberman/internal/adapters/accounts"
	"github.com/flowck/doberman/internal/adapters/events"
	"github.com/flowck/doberman/internal/adapters/monitors"
	"github.com/flowck/doberman/internal/adapters/resend"
	"github.com/flowck/doberman/internal/app"
	"github.com/flowck/doberman/internal/app/command"
	commontransaction "github.com/flowck/doberman/internal/common/adapters/transaction"
	"github.com/flowck/doberman/internal/common/ddd"
	commonhttpport "github.com/flowck/doberman/internal/common/httpport"
	"github.com/flowck/doberman/internal/common/logs"
	"github.com/flowck/doberman/internal/common/observability"
	"github.com/flowck/doberman/internal/common/psql"
	"github.com/flowck/doberman/internal/common/watermill_logger"
	httpport "github.com/flowck/doberman/internal/ports/http"

	"github.com/flowck/doberman/internal/ports/cron"
)

type Config struct {
	DatabaseUrl        string `envconfig:"DATABASE_URL"`
	AmqpURL            string `envconfig:"AMQP_URL"`
	Port               int    `envconfig:"PORT"`
	AllowedCorsOrigin  string `envconfig:"ALLOWED_CORS_ORIGIN"`
	MockResend         bool   `envconfig:"MOCK_RESEND"`
	ResendApiKey       string `envconfig:"RESEND_API_KEY"`
	NotificationsEmail string `envconfig:"NOTIFICATIONS_EMAIL"`
}

func main() {
	logger := logrus.New()
	wmLogger := watermill_logger.NewWatermillLogrusLogger(logger)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := &Config{}
	if err := envconfig.Process("", config); err != nil {
		logger.Fatalf("unable to read environment variables: %v", err)
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	amqpConfig := amqp.NewDurablePubSubConfig(config.AmqpURL, nil)

	publisher, err := amqp.NewPublisher(amqpConfig, wmLogger)
	if err != nil {
		logger.Fatalf("unable to create a amqp publisher: %v", err)
	}
	defer func() { _ = publisher.Close() }()

	db, err := psql.Connect(config.DatabaseUrl)
	if err != nil {
		logger.Fatal(err)
	}

	err = psql.ApplyMigrations(db, "/misc/sql/migrations")
	if err != nil {
		logger.Fatal(err)
	}

	monitorApplication := newApplication(db, publisher, logger, config)

	router := echo.New()
	httpPort := commonhttpport.NewPort(commonhttpport.PortConfig{
		Port:              8080,
		Logger:            logger,
		Ctx:               ctx,
		Router:            router,
		AllowedCorsOrigin: strings.Split(config.AllowedCorsOrigin, ","),
	})

	monitorRouter := router.Group("/monitor")
	httpport.RegisterHttpHandlers(monitorApplication, monitorRouter)

	go func() {
		logger.Infof("http port is running from http://localhost:%d", config.Port)
		if err = httpPort.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.WithError(err).Error("Something unexpected happened while running the http port")
		}
	}()

	deamon := cron.NewDaemon(cron.NewEnqueueMonitorsTask(monitorApplication), time.Second*10, logger)
	go func() {
		if err = deamon.Run(ctx); err != nil {
			logger.Error(err)
		}
	}()
	defer deamon.Stop()

	<-done
	terminationCtx, terminationCtxCancel := context.WithTimeout(ctx, time.Second*20)
	defer terminationCtxCancel()

	if err = httpPort.Stop(terminationCtx); err != nil {
		logger.WithError(err).Fatal("an error occurred while stopping the http port gracefully")
	}

	logger.Info("Exiting...")
}

func newApplication(db *sql.DB, publisher message.Publisher, logger *logs.Logger, config *Config) *app.App {
	eventPublisher := events.NewPublisher(publisher)
	monitorRepository := monitors.NewPsqlRepository(db)
	accountRepository := accounts.NewPsqlRepository(db)

	txProvider := commontransaction.NewPsqlProvider[command.TransactableAdapters](db, command.TransactableAdapters{
		EventPublisher:    eventPublisher,
		MonitorRepository: monitorRepository,
		AccountRepository: accountRepository,
	})

	resendClient := resendsdk.NewClient(config.ResendApiKey)
	notificationEmail, err := ddd.NewEmail(config.NotificationsEmail)
	if err != nil {
		panic(err)
	}

	application := &app.App{
		Queries: app.Queries{},
		Commands: app.Commands{
			CreateMonitor:   observability.NewCommandDecorator[command.CreateMonitor](command.NewCreateMonitorHandler(monitorRepository), logger, nil),
			EnqueueMonitors: observability.NewCommandDecorator[command.EnqueueMonitors](command.NewEnqueueMonitorsHandler(txProvider), logger, nil),

			// IAM
			CreateAccount:  observability.NewCommandDecorator[command.CreateAccount](command.NewCreateAccountHandler(txProvider, resend.NewService(resendClient, notificationEmail)), logger, nil),
			ConfirmAccount: observability.NewCommandDecorator[command.ConfirmAccount](command.NewConfirmAccountHandler(accountRepository), logger, nil),
		},
	}

	if config.MockResend {
		logger.Info("Using resend mock")
		application.Commands.CreateAccount = observability.NewCommandDecorator[command.CreateAccount](command.NewCreateAccountHandler(txProvider, resend.ServiceMock{}), logger, nil)
	}

	return application
}
