package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/flowck/dobermann/backend/internal/adapters/transaction"
	"github.com/flowck/dobermann/backend/internal/app"
	"github.com/flowck/dobermann/backend/internal/app/command"
	"github.com/flowck/dobermann/backend/internal/common/logs"
	"github.com/flowck/dobermann/backend/internal/common/observability"
	"github.com/flowck/dobermann/backend/internal/common/psql"
	httpport "github.com/flowck/dobermann/backend/internal/ports/http"
)

type Config struct {
	Port        int    `envconfig:"HTTP_PORT"`
	DebugMode   string `envconfig:"DEBUG_MODE"`
	DatabaseURL string `envconfig:"DATABASE_URL"`
}

func main() {
	config := &Config{}

	err := envconfig.Process("", config)
	if err != nil {
		panic(err)
	}

	logger := logs.New(config.DebugMode == "enabled")

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

	txProvider := transaction.NewPsqlProvider(db)

	application := &app.App{
		Commands: app.Commands{
			CreateAccount: observability.NewCommandDecorator[command.CreateAccount](command.NewCreateAccountHandler(txProvider), logger),
			CreateMonitor: nil,
		},
	}

	httpPort, err := httpport.NewPort(httpport.Config{
		Ctx:         ctx,
		Logger:      logger,
		Port:        config.Port,
		Application: application,
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

	<-done
	terminationCtx, terminationCancel := context.WithTimeout(context.Background(), time.Second*2)
	defer terminationCancel()

	err = httpPort.Stop(terminationCtx)
	if err != nil {
		logger.Fatalf("unable to gracefully shutdown the http port: %v", err)
	}
}
