package components_test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/flowck/dobermann/backend/tests/client"
)

var (
	// db  *sql.DB
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

func getClient(token string) *client.ClientWithResponses {
	host := fmt.Sprintf("http://localhost:%s", os.Getenv("HTTP_PORT"))
	newClient, err := client.NewClientWithResponses(host, client.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		if token != "" {
			req.Header.Add(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", token))
		}

		req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

		return nil
	}))

	if err != nil {
		log.Fatal(err)
	}

	return newClient
}

/*func unMarshallMessageToEvent[T any](m *message.Message) (T, error) {
	var event T
	err := json.Unmarshal(m.Payload, &event)
	if err != nil {
		return event, err
	}

	return event, nil
}*/
