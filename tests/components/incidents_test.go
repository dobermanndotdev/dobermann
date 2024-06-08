package components_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dobermanndotdev/dobermann/internal/adapters/models"
	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/tests"
	"github.com/dobermanndotdev/dobermann/tests/client"
)

func TestIncidents(t *testing.T) {
	user := createAccount(t)
	token := login(t, user.Email, user.Password)
	cli := getClient(token)

	monitor00Payload := fixtureMonitors(t, cli, 1)[0]
	monitor00 := getMonitorByEndpointUrl(t, monitor00Payload.EndpointUrl)
	monitor00ID, err := domain.NewIdFromString(monitor00.ID)
	require.NoError(t, err)

	t.Run("get_incident_by_id", func(t *testing.T) {
		incident00 := fixtureClient.FixtureAndInsertIncidents(t, monitor00ID, 1)[0]

		resp01, err := cli.GetIncidentByIDWithResponse(ctx, incident00.ID)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp01.StatusCode())

		foundIncident := resp01.JSON200.Data
		assert.Equal(t, incident00.ID, foundIncident.Id)
		assert.Equal(t, incident00.CreatedAt.UTC().Truncate(time.Second), foundIncident.CreatedAt.Truncate(time.Second))
		assert.Equal(t, incident00.Cause.String, foundIncident.Cause)
		assert.Equal(t, incident00.CheckedURL, foundIncident.CheckedUrl)
		assert.Equal(t, incident00.ResponseStatus.Int16, int16(*foundIncident.ResponseStatus))
	})

	t.Run("get_all_incidents", func(t *testing.T) {
		incident00 := fixtureClient.FixtureAndInsertIncidents(t, monitor00ID, 1)[0]
		resp01, err := cli.GetAllIncidentsWithResponse(ctx, &client.GetAllIncidentsParams{
			Page:  tests.ToPtr(1),
			Limit: tests.ToPtr(100),
		})
		require.NoError(t, err)

		matchedFixturedIncident := false
		for _, foundIncident := range resp01.JSON200.Data {
			if foundIncident.Id != incident00.ID {
				continue
			}

			assertIncident(t, incident00, foundIncident)
			matchedFixturedIncident = true
		}

		assert.True(t, matchedFixturedIncident)
	})
}

func assertIncident(t *testing.T, expected models.Incident, found client.Incident) {
	assert.Equal(t, expected.ID, found.Id)
	assert.Equal(t, expected.CreatedAt.UTC().Truncate(time.Second), found.CreatedAt.Truncate(time.Second))
	assert.Equal(t, expected.Cause.String, found.Cause)
	assert.Equal(t, expected.CheckedURL, found.CheckedUrl)
}
