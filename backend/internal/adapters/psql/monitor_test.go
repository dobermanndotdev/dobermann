package psql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/psql"
	"github.com/flowck/dobermann/backend/internal/app/query"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
	"github.com/flowck/dobermann/backend/tests"
)

func TestMonitorRepository_Lifecycle(t *testing.T) {
	const maxMonitors = 10
	repo := psql.NewMonitorRepository(db)
	acc := tests.FixtureAndInsertAccount(t, db)
	monitors := make(map[domain.ID]*monitor.Monitor)

	for i := 0; i < maxMonitors; i++ {
		m := tests.FixtureMonitor(t, acc.ID())
		monitors[m.ID()] = m
		require.NoError(t, repo.Insert(ctx, m))
	}

	limitPerPage := 2
	result, err := repo.FindAll(ctx, acc.ID(), query.PaginationParams{
		Page:  1,
		Limit: limitPerPage,
	})
	require.NoError(t, err)

	assert.Equal(t, int64(maxMonitors), result.TotalCount)
	assert.Equal(t, limitPerPage, result.PerPage)
	assert.Equal(t, maxMonitors/limitPerPage, result.PageCount)
	assert.Len(t, result.Data, limitPerPage)
}
