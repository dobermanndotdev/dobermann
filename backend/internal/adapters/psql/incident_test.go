package psql_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/psql"
	"github.com/flowck/dobermann/backend/tests"
)

func TestIncidentRepository_Lifecycle(t *testing.T) {
	monitorRepo := psql.NewMonitorRepository(db)
	incidentRepo := psql.NewIncidentRepository(db)

	account00 := tests.FixtureAndInsertAccount(t, db, true)
	monitor00 := tests.FixtureMonitor(t, account00)
	incident00 := tests.FixtureIncident(t)

	require.NoError(t, monitorRepo.Insert(ctx, monitor00))
	assert.NoError(t, incidentRepo.Create(ctx, monitor00.ID(), incident00))
}
