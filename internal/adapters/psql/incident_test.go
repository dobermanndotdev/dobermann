package psql_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dobermanndotdev/dobermann/internal/app/query"
	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
	"github.com/dobermanndotdev/dobermann/tests"
)

func TestIncidentRepository_Lifecycle(t *testing.T) {
	account00 := tests.FixtureAndInsertAccount(t, db, true)
	monitor00 := tests.FixtureMonitor(t, account00)
	incident00 := tests.FixtureIncident(t, monitor00.ID().String())

	require.NoError(t, monitorRepository.Insert(ctx, monitor00))
	assert.NoError(t, incidentRepository.Create(ctx, incident00))

	t.Run("find_incident_by_id", func(t *testing.T) {
		t.Parallel()

		found00, err := incidentRepository.FindByID(ctx, incident00.ID())
		require.NoError(t, err)
		assertIncident(t, incident00, found00)
	})

	t.Run("incident_not_found", func(t *testing.T) {
		t.Parallel()

		_, err := incidentRepository.FindByID(ctx, domain.NewID())
		assert.ErrorIs(t, err, monitor.ErrIncidentNotFound)
	})

	t.Run("all_incidents", func(t *testing.T) {
		t.Parallel()

		result, err := incidentRepository.FindAll(ctx, account00.ID(), query.PaginationParams{
			Page:  1,
			Limit: 100,
		})
		require.NoError(t, err)

		for _, foundIncident := range result.Data {
			require.Equal(t, monitor00.ID(), foundIncident.MonitorID())
		}
	})
}

func assertIncident(t *testing.T, expected, found *monitor.Incident) {
	t.Helper()

	assert.Equal(t, expected.ID(), found.ID())
	assert.Equal(t, expected.CheckedURL(), found.CheckedURL())
	assert.Equal(t, expected.Cause(), found.Cause())

	//POST_IDEA?: comparing dates in tests
	assert.Equal(t, expected.CreatedAt().Truncate(time.Second), found.CreatedAt().Truncate(time.Second))
}
