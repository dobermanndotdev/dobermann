package components_test

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMonitors(t *testing.T) {
	user := createAccount(t)
	token := login(t, user.Email, user.Password)
	require.NotEmpty(t, token)
}
