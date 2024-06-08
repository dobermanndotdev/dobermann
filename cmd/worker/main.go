package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/dobermanndotdev/dobermann/internal/adapters/endpoint_checkers"
	"github.com/dobermanndotdev/dobermann/internal/adapters/events"
	"github.com/dobermanndotdev/dobermann/internal/adapters/psql"
	"github.com/dobermanndotdev/dobermann/internal/adapters/transaction"
	"github.com/dobermanndotdev/dobermann/internal/app"
	"github.com/dobermanndotdev/dobermann/internal/app/command"
	"github.com/dobermanndotdev/dobermann/internal/common/logs"
	"github.com/dobermanndotdev/dobermann/internal/common/messaging"
	"github.com/dobermanndotdev/dobermann/internal/common/observability"
	"github.com/dobermanndotdev/dobermann/internal/common/postgres"
	"github.com/dobermanndotdev/dobermann/internal/ports/cron"
)

var Version = "development"

type Config struct {
	AmqpUrl                       string `envconfig:"AMQP_URL"`
	Port                          int    `envconfig:"HTTP_PORT"`
	DebugMode                     string `envconfig:"DEBUG_MODE"`
	DatabaseURL                   string `envconfig:"DATABASE_URL"`
	Region                        string `envconfig:"WORKER_REGION" required:"true"`
	IsProductionMode              bool   `envconfig:"PRODUCTION_MODE"`
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
	logger.Infof("Worker is running from region %s", config.Region)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	publisher, err := messaging.NewAmqpPublisher(config.AmqpUrl, logger)
	if err != nil {
		logger.Fatal(err)
	}

	db, err := postgres.Connect(config.DatabaseURL)
	if err != nil {
		logger.Fatal(err)
	}

	err = postgres.ApplyMigrations(db, "misc/sql/migrations")
	if err != nil {
		logger.Fatal(err)
	}

	httpChecker, err := endpoint_checkers.NewHttpChecker(config.Region, config.EndpointCheckTimeoutInSeconds, logger)
	if err != nil {
		logger.Fatal(err)
	}

	monitorRepository := psql.NewMonitorRepository(db)
	eventPublisher := events.NewPublisher(publisher)
	txProvider := transaction.NewPsqlProvider(db, publisher, logger)

	application := &app.App{
		Commands: app.Commands{
			BulkCheckEndpoints: observability.NewCommandDecorator[command.BulkCheckEndpoints](command.NewBulkCheckEndpointsHandler(httpChecker, txProvider, eventPublisher, monitorRepository), logger),
		},
	}

	cronService := cron.NewService(application, config.Region, config.IsProductionMode)

	go func() {
		err = cronService.Start(ctx)
		if err != nil {
			logger.Errorf("cron service stopped: %v", err)
		}
	}()

	logger.WithFields(logs.Fields{
		"version": Version,
		"region":  config.Region,
	}).Info("The service is running")

	<-done
	err = cronService.Stop()
	if err != nil {
		logger.Fatalf("unable to gracefully stop the cron service: %v", err)
	}

	logger.Info("The worker has been terminated gracefully")
}
