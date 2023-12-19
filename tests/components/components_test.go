package components_test

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq"

	"github.com/flowck/dobermann/backend/internal/common/postgres"
	"github.com/flowck/dobermann/backend/tests"
	"github.com/flowck/dobermann/backend/tests/client"
)

var (
	db            *sql.DB
	ctx           context.Context
	fixtureClient tests.FixtureClient
)

func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	var err error
	db, err = postgres.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	err = postgres.ApplyMigrations(db, "../../misc/sql/migrations")
	if err != nil {
		panic(err)
	}

	fixtureClient = tests.FixtureClient{
		Db:  db,
		Ctx: ctx,
	}

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
