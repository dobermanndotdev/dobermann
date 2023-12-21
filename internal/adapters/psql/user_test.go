package psql_test

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/psql"
	"github.com/flowck/dobermann/backend/internal/domain/account"
	"github.com/flowck/dobermann/backend/tests"
)

func TestUserRepository_Lifecycle(t *testing.T) {
	acc := tests.FixtureAndInsertAccount(t, db, false)
	user, err := acc.FirstAccountOwner()
	require.NoError(t, err)

	repo := psql.NewUserRepository(db)
	require.NoError(t, repo.Insert(ctx, user))
	require.ErrorIs(t, repo.Insert(ctx, user), account.ErrAccountExists)

	t.Run("find_by_email", func(t *testing.T) {
		userFound, err := repo.FindByEmail(ctx, user.Email())
		require.NoError(t, err)
		assert.Equal(t, user.ID(), userFound.ID())
	})

	t.Run("find_by_email_error_user_not_found", func(t *testing.T) {
		fakeEmail, err := account.NewEmail(gofakeit.Email())
		require.NoError(t, err)

		_, err = repo.FindByEmail(ctx, fakeEmail)
		assert.ErrorIs(t, err, account.ErrUserNotFound)
	})
}
