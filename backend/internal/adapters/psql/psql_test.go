package psql_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/flowck/dobermann/backend/internal/common/postgres"
)

var (
	db  *sql.DB
	ctx context.Context
)

// Set up file
func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	var err error
	db, err = postgres.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
