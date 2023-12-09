package main

import (
	"context"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/flowck/dobermann/backend/internal/adapters/endpoint_checkers"
	"github.com/flowck/dobermann/backend/internal/adapters/transaction"
	"github.com/flowck/dobermann/backend/internal/app"
	"github.com/flowck/dobermann/backend/internal/app/command"
	"github.com/flowck/dobermann/backend/internal/common/logs"
	"github.com/flowck/dobermann/backend/internal/common/messaging"
	"github.com/flowck/dobermann/backend/internal/common/observability"
	"github.com/flowck/dobermann/backend/internal/common/postgres"
	"github.com/flowck/dobermann/backend/internal/ports/cron"
)

type Config struct {
	AmqpUrl     string `envconfig:"AMQP_URL"`
	Port        int    `envconfig:"HTTP_PORT"`
	DebugMode   string `envconfig:"DEBUG_MODE"`
	DatabaseURL string `envconfig:"DATABASE_URL"`
	Region      string `envconfig:"FLY_REGION" required:"true"`
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

	httpChecker := endpoint_checkers.NewHttpChecker()
	txProvider := transaction.NewPsqlProvider(db, publisher)

	application := &app.App{
		Commands: app.Commands{
			BulkCheckEndpoints: observability.NewCommandDecorator[command.BulkCheckEndpoints](command.NewBulkCheckEndpointsHandler(httpChecker, txProvider), logger),
		},
	}

	cronService := cron.NewService(application, config.Region)

	go func() {
		err = cronService.Start(ctx)
		if err != nil {
			logger.Errorf("cron service stopped: %v", err)
		}
	}()

	<-done
	err = cronService.Stop()
	if err != nil {
		logger.Fatalf("unable to gracefully stop the cron service: %v", err)
	}

	logger.Info("The worker has been terminated gracefully")
}
