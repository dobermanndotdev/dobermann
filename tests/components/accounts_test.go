package components_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dobermanndotdev/dobermann/internal/adapters/models"
)

func TestAccounts_Lifecycle(t *testing.T) {
	userPayload := createAccount(t)
	token := login(t, userPayload.Email, userPayload.Password)
	cli := getClient(token)

	t.Run("get_profile_details", func(t *testing.T) {
		resp, err := cli.GetProfileDetailsWithResponse(ctx)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode())

		require.NotEmpty(t, resp.JSON200.Data)
		user := getUserByEmail(t, userPayload.Email)

		assert.Equal(t, user.ID, resp.JSON200.Data.Id)
		assert.Equal(t, user.FirstName.String, resp.JSON200.Data.FirstName)
		assert.Equal(t, user.LastName.String, resp.JSON200.Data.LastName)
		assert.Equal(t, user.Role, resp.JSON200.Data.Role)
		assert.Equal(t, user.CreatedAt.UTC().Truncate(time.Second), resp.JSON200.Data.CreatedAt.Truncate(time.Second))
	})
}

func getUserByEmail(t *testing.T, email string) *models.User {
	model, err := models.Users(models.UserWhere.Email.EQ(email)).One(ctx, db)
	require.NoError(t, err)

	return model
}
