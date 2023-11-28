package monitors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/monitors"
	"github.com/flowck/dobermann/backend/tests"
)

func TestIncidentRepository_Lifecycle(t *testing.T) {
	monitorRepo := monitors.NewPsqlRepository(db)
	incidentRepo := monitors.NewIncidentRepository(db)

	acc := tests.FixtureAndInsertAccount(t, db)
	monitor00 := tests.FixtureMonitor(t, acc.ID())
	incident00 := tests.FixtureIncident(t)

	require.NoError(t, monitorRepo.Insert(ctx, monitor00))
	assert.NoError(t, incidentRepo.Create(ctx, monitor00.ID(), incident00))
}
