package components_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
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
			EndpointUrl:            endpointUrl,
			CheckIntervalInSeconds: 180,
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp01.StatusCode)
		assert.Eventually(t, assertMonitorHasBeenChecked(t, endpointUrl), time.Second*5, time.Millisecond*250)

		resp02, err := cli.CreateMonitor(ctx, client.CreateMonitorRequest{
			EndpointUrl:            endpointUrl,
			CheckIntervalInSeconds: 29,
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp02.StatusCode)
	})

	t.Run("create_monitor_with_and_endpoint_down", func(t *testing.T) {
		t.Parallel()
		endpointUrl := endpointUrlGenerator(false)
		resp01, err := cli.CreateMonitor(ctx, client.CreateMonitorRequest{
			EndpointUrl:            endpointUrl,
			CheckIntervalInSeconds: 30,
		})
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp01.StatusCode)

		//POST_IDEA: How to test eventually consistent features
		createdMonitor := getMonitorByEndpointUrl(t, endpointUrl)
		assert.Eventually(t, assertIncidentHasBeenCreated(createdMonitor.ID), time.Second*10, time.Millisecond*250)

		//TODO: extend this test by refactoring the simulator to be more generic such that it can be used as a mock service
		//TODO: for all external services such as Resend and then used to check whether a certain call has been made
		//TODO: expected an incident action to have been created
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
		monitorPayload := fixtureMonitors(t, cli, 1)[0]
		monitor00 := getMonitorByEndpointUrl(t, monitorPayload.EndpointUrl)

		resp01, err := cli.GetMonitorByIDWithResponse(ctx, monitor00.ID)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp01.StatusCode())

		assert.Equal(t, monitor00.ID, resp01.JSON200.Data.Id)
		assert.Equal(t, monitor00.EndpointURL, resp01.JSON200.Data.EndpointUrl)
		assert.False(t, monitor00.IsPaused)
		assert.Equal(t, monitorPayload.CheckIntervalInSeconds, monitor00.CheckIntervalInSeconds)
	})

	t.Run("pause_and_unpause_monitor", func(t *testing.T) {
		monitorPayload := fixtureMonitors(t, cli, 1)[0]
		monitor00 := getMonitorByEndpointUrl(t, monitorPayload.EndpointUrl)
		require.False(t, monitor00.IsPaused)

		resp01, err := cli.ToggleMonitorPause(ctx, monitor00.ID, client.ToggleMonitorPauseRequest{Pause: true})
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, resp01.StatusCode)

		monitor00 = getMonitorByEndpointUrl(t, monitorPayload.EndpointUrl)
		assert.True(t, monitor00.IsPaused)

		resp02, err := cli.ToggleMonitorPause(ctx, monitor00.ID, client.ToggleMonitorPauseRequest{Pause: false})
		require.NoError(t, err)
		require.Equal(t, http.StatusNoContent, resp02.StatusCode)

		monitor00 = getMonitorByEndpointUrl(t, monitorPayload.EndpointUrl)
		assert.False(t, monitor00.IsPaused)
	})

	t.Run("delete_monitor_by_id", func(t *testing.T) {
		monitorPayload := fixtureMonitors(t, cli, 1)[0]
		monitor00 := getMonitorByEndpointUrl(t, monitorPayload.EndpointUrl)
		resp, err := cli.DeleteMonitor(ctx, monitor00.ID)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
	})

	t.Run("edit_monitor", func(t *testing.T) {
		monitorPayload := fixtureMonitors(t, cli, 1)[0]
		monitor00 := getMonitorByEndpointUrl(t, monitorPayload.EndpointUrl)

		newDetails := client.EditMonitorRequest{
			CheckIntervalInSeconds: 300,
			EndpointUrl:            endpointUrlGenerator(true),
		}

		resp, err := cli.EditMonitor(ctx, monitor00.ID, newDetails)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		monitor00 = getMonitorByEndpointUrl(t, newDetails.EndpointUrl)
		assert.Equal(t, newDetails.EndpointUrl, monitor00.EndpointURL)
		assert.Equal(t, newDetails.CheckIntervalInSeconds, monitor00.CheckIntervalInSeconds)
	})

	t.Run("get_monitor_response_time_stats", func(t *testing.T) {
		monitorPayload := fixtureMonitors(t, cli, 1)[0]
		monitor00 := getMonitorByEndpointUrl(t, monitorPayload.EndpointUrl)
		monitorID, err := domain.NewIdFromString(monitor00.ID)
		require.NoError(t, err)

		fixtureClient.FixtureCheckResults(t, monitorID, time.Now().Add(time.Hour*-24*7), 7, 2, 30)

		resp01, err := cli.GetMonitorResponseTimeStatsWithResponse(ctx, monitor00.ID, &client.GetMonitorResponseTimeStatsParams{
			RangeInDays: tests.ToPtr(7),
		})
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp01.StatusCode())
		require.NotEmpty(t, resp01.JSON200.Data)
	})
}

func endpointUrlGenerator(isUp bool) string {
	isUpParam := "false"

	if isUp {
		isUpParam = "true"
	}

	return fmt.Sprintf("%s?id=%s&is_up=%s", tests.SimulatorEndpointUrl, domain.NewID().String(), isUpParam)
}

func fixtureMonitors(t *testing.T, cli *client.ClientWithResponses, maxEndpoints int) []client.CreateMonitorRequest {
	endpoints := make([]client.CreateMonitorRequest, maxEndpoints)
	checkIntervalInSeconds := []int{30, 60, 180, 300, 900, 1800, 3600}

	var endpointUrl string
	for i := 0; i < maxEndpoints; i++ {
		endpointUrl = fmt.Sprintf("%s?id=%s", tests.SimulatorEndpointUrl, domain.NewID().String())
		endpoints[i] = client.CreateMonitorRequest{
			EndpointUrl:            endpointUrl,
			CheckIntervalInSeconds: checkIntervalInSeconds[gofakeit.Number(0, len(checkIntervalInSeconds)-1)],
		}

		resp01, err := cli.CreateMonitor(ctx, endpoints[i])
		require.NoError(t, err)
		require.Equal(t, http.StatusCreated, resp01.StatusCode)
	}

	return endpoints
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
