package components_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill-amqp/v2/pkg/amqp"
	"github.com/ThreeDotsLabs/watermill/message"
	_ "github.com/lib/pq"

	"github.com/flowck/doberman/internal/common/logs"
	"github.com/flowck/doberman/internal/common/psql"
	"github.com/flowck/doberman/internal/common/watermill_logger"
	"github.com/flowck/doberman/tests/client"
)

var (
	db         *sql.DB
	ctx        context.Context
	cli        *client.ClientWithResponses
	subscriber message.Subscriber
	logger     = logs.New(true)
	wmLogger   = watermill_logger.NewWatermillLogrusLogger(logger)
)

func TestMain(m *testing.M) {
	var err error
	host := "http://localhost:8080"

	cli, err = client.NewClientWithResponses(fmt.Sprintf("%s/monitor", host))
	if err != nil {
		panic(err)
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	subscriber, err = amqp.NewSubscriber(
		amqp.NewDurablePubSubConfig(
			os.Getenv("AMQP_URL"),
			amqp.GenerateQueueNameTopicNameWithSuffix("worker_url_crawlers"),
		),
		wmLogger,
	)
	if err != nil {
		panic(err)
	}

	db, err = psql.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	err = psql.ApplyMigrations(db, "../../misc/sql/migrations")
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}

func unMarshallMessageToEvent[T any](m *message.Message) (T, error) {
	var event T
	err := json.Unmarshal(m.Payload, &event)
	if err != nil {
		return event, err
	}

	return event, nil
}
