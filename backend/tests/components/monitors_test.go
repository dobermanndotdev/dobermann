package components_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/tests"
	"github.com/flowck/dobermann/backend/tests/client"
)

func TestMonitors(t *testing.T) {
	user := createAccount(t)
	token := login(t, user.Email, user.Password)
	cli := getClient(token)

	endpointUrl := fmt.Sprintf("%s#id=%s", tests.SimulatorEndpointUrl, domain.NewID().String())
	resp01, err := cli.CreateMonitor(ctx, client.CreateMonitorRequest{
		EndpointUrl: endpointUrl,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp01.StatusCode)
	assert.Eventually(t, assertMonitorIsWillBeChecked(t, endpointUrl), time.Second*5, time.Millisecond*250)
}

func assertMonitorIsWillBeChecked(t *testing.T, endpointUrl string) func() bool {
	return func() bool {
		model, err := models.Monitors(models.MonitorWhere.EndpointURL.EQ(endpointUrl)).One(ctx, db)
		require.NoError(t, err)

		return model.LastCheckedAt.Ptr() != nil
	}
}
