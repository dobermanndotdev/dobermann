package psql_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/psql"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
	"github.com/flowck/dobermann/backend/tests"
)

func TestIncidentRepository_Lifecycle(t *testing.T) {
	monitorRepo := psql.NewMonitorRepository(db)

	account00 := tests.FixtureAndInsertAccount(t, db, true)
	monitor00 := tests.FixtureMonitor(t, account00)
	incident00 := tests.FixtureIncident(t)

	require.NoError(t, monitorRepo.Insert(ctx, monitor00))
	assert.NoError(t, incidentRepository.Create(ctx, monitor00.ID(), incident00))

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
}

func assertIncident(t *testing.T, expected, found *monitor.Incident) {
	t.Helper()

	assert.Equal(t, expected.ID(), found.ID())
	assert.Equal(t, expected.CheckedURL(), found.CheckedURL())
	assert.Equal(t, expected.Details().Cause, found.Details().Cause)

	//POST_IDEA?: comparing dates in tests
	assert.Equal(t, expected.CreatedAt().Truncate(time.Second), found.CreatedAt().Truncate(time.Second))
}
