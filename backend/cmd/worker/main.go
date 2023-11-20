package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"

	"github.com/flowck/doberman/internal/adapters/monitors"
	"github.com/flowck/doberman/internal/app"
	monitorcommand "github.com/flowck/doberman/internal/app/command"
	commontransaction "github.com/flowck/doberman/internal/common/adapters/transaction"
	"github.com/flowck/doberman/internal/common/logs"
	"github.com/flowck/doberman/internal/common/psql"
	"github.com/flowck/doberman/internal/common/watermill_logger"
	amqpport "github.com/flowck/doberman/internal/ports/amqp"
)

const (
	consumerGroupUrlCrawlers = "worker_url_crawlers"
)

type Config struct {
	DatabaseUrl string `envconfig:"DATABASE_URL"`
	AmqpURL     string `envconfig:"AMQP_URL"`
}

func main() {
	logger := logs.New(true)
	wmLogger := watermill_logger.NewWatermillLogrusLogger(logger)

	config := &Config{}
	if err := envconfig.Process("", config); err != nil {
		logger.Fatalf("unable to read environment variables: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	db, err := psql.Connect(config.DatabaseUrl)
	if err != nil {
		logger.Fatal(err)
	}

	monitorRepository := monitors.NewPsqlRepository(db)

	txProvider := commontransaction.NewPsqlProvider[monitorcommand.TransactableAdapters](db, monitorcommand.TransactableAdapters{
		EventPublisher:    nil,
		MonitorRepository: monitorRepository,
	})

	logger.Info(txProvider)

	application := &app.App{
		Queries:  app.Queries{},
		Commands: app.Commands{},
	}

	router, err := message.NewRouter(message.RouterConfig{}, wmLogger)
	if err != nil {
		panic(err)
	}

	publisher, err := amqp.NewPublisher(
		amqp.NewDurableQueueConfig(config.AmqpURL),
		wmLogger,
	)
	if err != nil {
		logger.Fatalf("unable to create publisher: %v", err)
	}

	poisonQueueMiddleware, err := middleware.PoisonQueue(publisher, "failed_crawling")
	if err != nil {
		logger.Fatalf("unable to create poison queue: %v", err)
	}

	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		middleware.CorrelationID,
		middleware.Retry{
			MaxRetries:      3,
			InitialInterval: time.Second * 5,
			Logger:          wmLogger,
		}.Middleware,
		poisonQueueMiddleware,
		middleware.NewThrottle(10, time.Second).Middleware,
		middleware.Recoverer,
	)

	subscriber, err := amqp.NewSubscriber(
		amqp.NewDurablePubSubConfig(
			config.AmqpURL,
			amqp.GenerateQueueNameTopicNameWithSuffix(consumerGroupUrlCrawlers),
		),
		wmLogger,
	)
	if err != nil {
		panic(err)
	}

	for _, handler := range amqpport.NewHandlers(application) {
		router.AddNoPublisherHandler(handler.Name(), handler.EventName(), subscriber, handler.Handle)
	}

	if err = router.Run(ctx); err != nil {
		logger.WithError(err).Error("An error occurred while running watermill router")
	}

	<-done
	logger.Info("Exiting...")
}
