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
		endpointUrl := endpointUrlGenerator(false)
		resp01, err := cli.CreateMonitor(ctx, client.CreateMonitorRequest{
			EndpointUrl: endpointUrl,
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp01.StatusCode)
		assert.Eventually(t, assertMonitorHasBeenChecked(t, endpointUrl), time.Second*5, time.Millisecond*250)
	})

	t.Run("create_monitor_with_and_endpoint_down", func(t *testing.T) {
		t.Parallel()
		endpointUrl := endpointUrlGenerator(false)
		resp01, err := cli.CreateMonitor(ctx, client.CreateMonitorRequest{
			EndpointUrl: endpointUrl,
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp01.StatusCode)

		//POST_IDEA: How to test eventually consistent features
		createdMonitor := getMonitorByEndpointUrl(t, endpointUrl)
		assert.Eventually(t, assertIncidentHasBeenCreated(createdMonitor.ID), time.Second*10, time.Millisecond*250)
	})

	t.Run("get_all_monitors", func(t *testing.T) {
		// not parallel

		fixtureMonitors(t, cli, 5)

		resp01, err := cli.GetAllMonitorsWithResponse(ctx, &client.GetAllMonitorsParams{
			Page:  tests.ToPtr(1),
			Limit: tests.ToPtr(100),
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp01.StatusCode())
		require.NotEmpty(t, resp01.JSON200.Data)

		for _, m := range resp01.JSON200.Data {
			assert.NotEmpty(t, m.Id)
			assert.NotEmpty(t, m.EndpointUrl)
			assert.NotEmpty(t, m.CreatedAt)
		}

		resp02, err := cli.GetAllMonitorsWithResponse(ctx, &client.GetAllMonitorsParams{
			Page:  tests.ToPtr(1),
			Limit: tests.ToPtr(10),
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp02.StatusCode())

		require.Len(t, resp02.JSON200.Data, 5)
	})

	t.Run("get_monitor_by_id", func(t *testing.T) {
		endpointUrl := fixtureMonitors(t, cli, 1)[0]
		monitor00 := getMonitorByEndpointUrl(t, endpointUrl)

		resp01, err := cli.GetMonitorByIDWithResponse(ctx, monitor00.ID)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp01.StatusCode())

		assert.Equal(t, monitor00.ID, resp01.JSON200.Data.Id)
		assert.Equal(t, monitor00.EndpointURL, resp01.JSON200.Data.EndpointUrl)
	})
}

func endpointUrlGenerator(isUp bool) string {
	isUpParam := "false"

	if isUp {
		isUpParam = "true"
	}

	return fmt.Sprintf("%s?id=%s&is_up=%s", tests.SimulatorEndpointUrl, domain.NewID().String(), isUpParam)
}

func fixtureMonitors(t *testing.T, cli *client.ClientWithResponses, maxEndpoints int) []string {
	endpointUrls := make([]string, maxEndpoints)

	var endpointUrl string
	for i := 0; i < maxEndpoints; i++ {
		endpointUrl = fmt.Sprintf("%s?id=%s", tests.SimulatorEndpointUrl, domain.NewID().String())
		endpointUrls[i] = endpointUrl

		resp01, err := cli.CreateMonitor(ctx, client.CreateMonitorRequest{
			EndpointUrl: endpointUrl,
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp01.StatusCode)
	}

	return endpointUrls
}

func assertMonitorHasBeenChecked(t *testing.T, endpointUrl string) func() bool {
	return func() bool {
		model, err := models.Monitors(models.MonitorWhere.EndpointURL.EQ(endpointUrl)).One(ctx, db)
		require.NoError(t, err)

		return model.LastCheckedAt.Ptr() != nil
	}
}

func assertIncidentHasBeenCreated(monitorID string) func() bool {
	return func() bool {
		_, err := models.Incidents(
			models.IncidentWhere.MonitorID.EQ(monitorID),
			qm.OrderBy("created_at DESC"),
		).One(ctx, db)

		return err == nil
	}
}

func getMonitorByEndpointUrl(t *testing.T, endpointUrl string) *models.Monitor {
	model, err := models.Monitors(models.MonitorWhere.EndpointURL.EQ(endpointUrl)).One(ctx, db)
	require.NoError(t, err)
	return model
}
