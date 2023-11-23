package components_test

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ThreeDotsLabs/watermill/message"
	_ "github.com/lib/pq"

	"github.com/flowck/doberman/tests/client"
)

var (
	db  *sql.DB
	ctx context.Context
	cli *client.ClientWithResponses
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
