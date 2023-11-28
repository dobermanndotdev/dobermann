package monitors_test

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/monitors"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
	"github.com/flowck/dobermann/backend/tests"
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
