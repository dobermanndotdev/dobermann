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
	account00 := tests.FixtureAndInsertAccount(t, db, true)
	monitors := make(map[domain.ID]*monitor.Monitor)

	for i := 0; i < maxMonitors; i++ {
		m := tests.FixtureMonitor(t, account00)
		monitors[m.ID()] = m
		require.NoError(t, repo.Insert(ctx, m))
	}

	limitPerPage := 2
	result, err := repo.FindAll(ctx, account00.ID(), query.PaginationParams{
		Page:  1,
		Limit: limitPerPage,
	})
	require.NoError(t, err)

	assert.Equal(t, int64(maxMonitors), result.TotalCount)
	assert.Equal(t, limitPerPage, result.PerPage)
	assert.Equal(t, maxMonitors/limitPerPage, result.PageCount)
	assert.Len(t, result.Data, limitPerPage)

	t.Run("find_by_id", func(t *testing.T) {
		expected := result.Data[0]
		var found *monitor.Monitor

		found, err = repo.FindByID(ctx, expected.ID())
		require.NoError(t, err)

		assertMonitor(t, expected, found)

		owner, err := account00.FirstAccountOwner()
		require.NoError(t, err)

		assert.Equal(t, found.Subscribers()[0].UserID(), owner.ID())
	})

	t.Run("error_not_found_while_finding_by_id", func(t *testing.T) {
		_, err = repo.FindByID(ctx, domain.NewID())
		assert.ErrorIs(t, err, monitor.ErrMonitorNotFound)
	})
}

func assertMonitor(t *testing.T, expected, found *monitor.Monitor) {
	t.Helper()

	assert.Equal(t, expected.ID(), found.ID())
	assert.Equal(t, expected.AccountID(), found.AccountID())
	assert.Equal(t, expected.EndpointUrl(), found.EndpointUrl())
	assert.Equal(t, expected.LastCheckedAt(), found.LastCheckedAt())
	assert.Equal(t, expected.CreatedAt(), found.CreatedAt())
	assert.Equal(t, expected.IsEndpointUp(), found.IsEndpointUp())
	assert.Equal(t, expected.CheckInterval(), found.CheckInterval())
	assert.NotEmpty(t, found.Subscribers(), "has subscribers")
}
