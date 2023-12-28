package components_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/domain"
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
		assert.Equal(t, int(incident00.ResponseStatus), foundIncident.ResponseStatus)
		assert.Equal(t, incident00.RequestHeaders.String, foundIncident.RequestHeaders)
		assert.Equal(t, incident00.ResponseHeaders.String, foundIncident.ResponseHeaders)
	})
}
