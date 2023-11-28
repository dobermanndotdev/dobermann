package components_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/tests"
	"github.com/flowck/dobermann/backend/tests/client"
)

func TestMonitors(t *testing.T) {
	user := createAccount(t)
	token := login(t, user.Email, user.Password)
	cli := getClient(token)

	t.Run("create_monitor", func(t *testing.T) {
		t.Parallel()
		endpointUrl := fmt.Sprintf("%s#id=%s", tests.SimulatorEndpointUrl, domain.NewID().String())
		resp01, err := cli.CreateMonitor(ctx, client.CreateMonitorRequest{
			EndpointUrl: endpointUrl,
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp01.StatusCode)
		assert.Eventually(t, assertMonitorHasBeenChecked(t, endpointUrl), time.Second*5, time.Millisecond*250)
	})

	t.Run("create_monitor_with_and_endpoint_down", func(t *testing.T) {
		t.Parallel()
		endpointUrl := fmt.Sprintf("%s#id=%s&is_up=false", tests.SimulatorEndpointUrl, domain.NewID().String())
		resp01, err := cli.CreateMonitor(ctx, client.CreateMonitorRequest{
			EndpointUrl: endpointUrl,
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp01.StatusCode)
		assert.Eventually(t, assertIncidentExists(t, endpointUrl), time.Second*10, time.Millisecond*250)
	})
}

func assertMonitorHasBeenChecked(t *testing.T, endpointUrl string) func() bool {
	return func() bool {
		model, err := models.Monitors(models.MonitorWhere.EndpointURL.EQ(endpointUrl)).One(ctx, db)
		require.NoError(t, err)

		return model.LastCheckedAt.Ptr() != nil
	}
}

func assertIncidentExists(t *testing.T, endpointUrl string) func() bool {
	return func() bool {
		model, err := models.Monitors(models.MonitorWhere.EndpointURL.EQ(endpointUrl)).One(ctx, db)
		require.NoError(t, err)

		_, err = models.Incidents(
			models.IncidentWhere.MonitorID.EQ(model.ID),
			qm.OrderBy("created_at DESC"),
		).One(ctx, db)

		return err == nil
	}
}
