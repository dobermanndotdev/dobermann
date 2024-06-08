package components_test

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dobermanndotdev/dobermann/tests"
	"github.com/dobermanndotdev/dobermann/tests/client"
)

func TestAccessToProtectedEndpoints(t *testing.T) {
	resp01, err := getClient("").CreateMonitor(ctx, client.CreateMonitorRequest{
		EndpointUrl:            tests.SimulatorEndpointUrl,
		CheckIntervalInSeconds: 30,
	})
	require.NoError(t, err)
	assert.Equal(t, http.StatusForbidden, resp01.StatusCode)

	acc := createAccount(t)
	token := login(t, acc.Email, acc.Password)

	resp02, err := getClient(token).CreateMonitor(ctx, client.CreateMonitorRequest{
		EndpointUrl:            fmt.Sprintf("%s#testAcccessToProtectedEndpoints", tests.SimulatorEndpointUrl),
		CheckIntervalInSeconds: 30,
	})
	require.NoError(t, err)
	assert.NotEqual(t, http.StatusForbidden, resp02.StatusCode)
}

func createAccount(t *testing.T) client.CreateAccountRequest {
	payload := client.CreateAccountRequest{
		Email:       gofakeit.Email(),
		Password:    tests.FixturePassword(),
		AccountName: gofakeit.Company(),
	}

	resp, err := getClient("").CreateAccount(ctx, payload)
	require.NoError(t, err)
	require.Equal(t, http.StatusCreated, resp.StatusCode)

	return payload
}

func login(t *testing.T, email, password string) string {
	resp, err := getClient("").LoginWithResponse(ctx, client.LogInRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, resp.StatusCode())
	require.NotEmpty(t, resp.JSON200.Token)

	return resp.JSON200.Token
}
