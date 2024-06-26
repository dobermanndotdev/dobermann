package psql_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/dobermanndotdev/dobermann/internal/adapters/models"
	"github.com/dobermanndotdev/dobermann/internal/app/query"
	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
	"github.com/dobermanndotdev/dobermann/tests"
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
	assertIncidentsAreSortedInDescOrder(t, result.Data[0].Incidents())

	t.Run("find_by_id", func(t *testing.T) {
		expected := result.Data[0]
		var found *monitor.Monitor

		incidents := make([]*monitor.Incident, 5)
		for i := 0; i < 5; i++ {
			incidents[i] = tests.FixtureIncident(t, expected.ID().String())
			if i != 4 {
				incidents[i].Resolve()
			}

			err = incidentRepository.Create(ctx, incidents[i])
			require.NoError(t, err)
		}

		found, err = monitorRepository.FindByID(ctx, expected.ID())
		require.NoError(t, err)

		assertMonitor(t, expected, found)

		owner, err := account00.FirstAccountOwner()
		require.NoError(t, err)
		assert.Equal(t, found.Subscribers()[0].UserID(), owner.ID())

		require.Len(t, found.Incidents(), 5)
		// assertIncident(t, incident00, found.Incidents()[0])
		assertIncidentsAreSortedInDescOrder(t, found.Incidents())
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

	err := monitorRepository.Delete(ctx, monitor00.ID())
	assert.NoError(t, err)
}

func TestMonitorRepository_ResponseTimeStats(t *testing.T) {
	account00 := tests.FixtureAndInsertAccount(t, db, true)
	now := time.Now()

	testCases := []struct {
		name                   string
		rangeInDays            int
		checksPerDay           int
		checkIntervalInSeconds int
		startCheckedAt         time.Time
		expectedCheckResults   int
	}{
		{
			name:                   "query_only_daily_check_results",
			rangeInDays:            1,
			checksPerDay:           48,
			checkIntervalInSeconds: 30,
			startCheckedAt:         time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, time.UTC),
			expectedCheckResults:   48,
		},
		{
			name:                   "query_weekly_check_results",
			rangeInDays:            7,
			checksPerDay:           8,
			checkIntervalInSeconds: 60 * 60, // 1h
			startCheckedAt:         time.Date(now.Year(), now.Month(), now.Day()-7, 0, 0, 0, 0, time.UTC),
			expectedCheckResults:   8 * 7,
		},
		{
			name:                   "query_monthly_check_results",
			rangeInDays:            31,
			checksPerDay:           8,
			checkIntervalInSeconds: 60 * 60, // 1h
			startCheckedAt:         time.Date(now.Year(), now.Month(), now.Day()-31, 0, 0, 0, 0, time.UTC),
			expectedCheckResults:   8 * 31,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			monitor00 := tests.FixtureMonitor(t, account00)
			require.NoError(t, monitorRepository.Insert(ctx, monitor00))
			fixtureClient.FixtureCheckResults(t, monitor00.ID(), tc.startCheckedAt, tc.rangeInDays, tc.checksPerDay, tc.checkIntervalInSeconds)

			result, err := monitorRepository.ResponseTimeStats(ctx, monitor00.ID(), tests.ToPtr(tc.rangeInDays))
			require.NoError(t, err)
			require.NotEmpty(t, result)

			assert.Equal(t, tc.expectedCheckResults, len(result))
		})
	}
}

func TestMonitorRepository_UpdateForCheck(t *testing.T) {
	accID := tests.FixtureAndInsertAccount(t, db, true).ID()
	monitor01 := fixtureMonitor(t, accID, time.Second*30, time.Now().Add(-time.Minute))

	err := monitorRepository.Insert(ctx, monitor01)
	require.NoError(t, err)

	saveIncident(t, monitor01.ID())

	var allFoundMonitors []*monitor.Monitor
	err = monitorRepository.UpdateForCheck(ctx, func(foundMonitors []*monitor.Monitor) error {
		allFoundMonitors = foundMonitors
		return nil
	})

	require.NotEmpty(t, allFoundMonitors)

	for _, m := range allFoundMonitors {
		assert.True(t, time.Since(*m.LastCheckedAt()) > m.CheckInterval(), "monitor's now() - lastCheckedAt must be equal or greater than check interval in seconds")

		if m.ID() == monitor01.ID() {
			assert.Len(t, m.Incidents(), 1)
		}
	}

	require.NoError(t, err)
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

func saveIncident(t *testing.T, monitorID domain.ID) {
	model := models.Incident{
		ID:             domain.NewID().String(),
		MonitorID:      monitorID.String(),
		ResolvedAt:     null.Time{},
		Cause:          null.String{},
		ResponseStatus: null.Int16From(http.StatusInternalServerError),
		CheckedURL:     gofakeit.URL(),
	}
	require.NoError(t, model.Insert(ctx, db, boil.Infer()))
}

func assertIncidentsAreSortedInDescOrder(t *testing.T, incidents []*monitor.Incident) {
	t.Helper()
	var prev *monitor.Incident

	for _, current := range incidents {
		if prev != nil {
			require.True(t, prev.CreatedAt().After(current.CreatedAt()))
		}
		prev = current
	}
}

func fixtureMonitor(
	t *testing.T,
	accID domain.ID,
	checkInterval time.Duration,
	lastCheckedAt time.Time,
) *monitor.Monitor {
	m, err := monitor.NewMonitor(
		domain.NewID(),
		tests.EndpointUrlGenerator(true),
		accID,
		true,
		false,
		nil,
		nil,
		time.Now(),
		checkInterval,
		tests.ToPtr(lastCheckedAt),
	)
	require.NoError(t, err)

	return m
}
