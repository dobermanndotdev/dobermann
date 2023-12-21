package psql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/app/query"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
	"github.com/flowck/dobermann/backend/tests"
)

func TestMonitorRepository_Lifecycle(t *testing.T) {
	const maxMonitors = 10
	account00 := tests.FixtureAndInsertAccount(t, db, true)
	monitors := make(map[domain.ID]*monitor.Monitor)

	for i := 0; i < maxMonitors; i++ {
		m := tests.FixtureMonitor(t, account00)
		monitors[m.ID()] = m
		require.NoError(t, monitorRepository.Insert(ctx, m))
	}

	limitPerPage := 2
	result, err := monitorRepository.FindAll(ctx, account00.ID(), query.PaginationParams{
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

		found, err = monitorRepository.FindByID(ctx, expected.ID())
		require.NoError(t, err)

		assertMonitor(t, expected, found)

		owner, err := account00.FirstAccountOwner()
		require.NoError(t, err)

		assert.Equal(t, found.Subscribers()[0].UserID(), owner.ID())
	})

	t.Run("error_not_found_while_finding_by_id", func(t *testing.T) {
		_, err = monitorRepository.FindByID(ctx, domain.NewID())
		assert.ErrorIs(t, err, monitor.ErrMonitorNotFound)
	})
}

func TestMonitorRepository_Delete(t *testing.T) {
	account00 := tests.FixtureAndInsertAccount(t, db, true)
	monitor00 := tests.FixtureMonitor(t, account00)
	require.NoError(t, monitorRepository.Insert(ctx, monitor00))

	fixtureClient.FixtureCheckResults(t, monitor00.ID(), 100, 2)
	fixtureClient.FixtureAndInsertIncidents(t, monitor00.ID(), 10)

	err := monitorRepository.Delete(ctx, monitor00.ID())
	assert.NoError(t, err)

	_, err = monitorRepository.FindByID(ctx, monitor00.ID())
	assert.ErrorIs(t, err, monitor.ErrMonitorNotFound, "monitor must be deleted")

	incidentRows, err := models.Incidents(models.IncidentWhere.MonitorID.EQ(monitor00.ID().String())).All(ctx, db)
	require.NoError(t, err)
	assert.Nil(t, incidentRows, "incidents must be deleted")

	checkResultRows, err := models.MonitorCheckResults(
		models.MonitorCheckResultWhere.MonitorID.EQ(monitor00.ID().String()),
	).All(ctx, db)
	require.NoError(t, err)
	assert.Nil(t, checkResultRows, "check results must be deleted")
}

func TestMonitorRepository_ResponseTimeStats(t *testing.T) {
	account00 := tests.FixtureAndInsertAccount(t, db, true)
	monitor00 := tests.FixtureMonitor(t, account00)
	require.NoError(t, monitorRepository.Insert(ctx, monitor00))
	expectedAvgResponseTimeInMs := int16(200)
	rangeInDays := 10

	fixtureClient.FixtureCheckResults(t, monitor00.ID(), expectedAvgResponseTimeInMs, rangeInDays)

	responseTimeStats, err := monitorRepository.ResponseTimeStats(ctx, query.ResponseTimeStatsOptions{
		RangeInDays: rangeInDays,
		MonitorID:   monitor00.ID(),
	})
	require.NoError(t, err)
	require.NotEmpty(t, responseTimeStats.ResponseTimePerRegion)

	for _, region := range responseTimeStats.ResponseTimePerRegion {
		require.Len(t, region.Data, rangeInDays)

		for _, avgResponseTimePerDay := range region.Data {
			assert.Equal(t, expectedAvgResponseTimeInMs, avgResponseTimePerDay.Value)
		}
	}
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
