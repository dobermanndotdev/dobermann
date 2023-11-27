package monitors_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/monitors"
	"github.com/flowck/dobermann/backend/internal/common/psql"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
	"github.com/flowck/dobermann/backend/tests"
)

var (
	db  *sql.DB
	ctx context.Context
)

func TestPsqlRepository_Lifecycle(t *testing.T) {
	repo := monitors.NewPsqlRepository(db)

	acc := tests.FixtureAndInsertAccount(t, db)
	monitorFixture := tests.FixtureMonitor(t, acc.ID())
	require.NoError(t, repo.Insert(ctx, monitorFixture))

	err := repo.Update(ctx, monitorFixture.ID(), func(foundMonitor *monitor.Monitor) error {
		foundMonitor.SetEndpointCheckResult(true)
		return nil
	})
	require.NoError(t, err)
}

func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	var err error
	db, err = psql.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
