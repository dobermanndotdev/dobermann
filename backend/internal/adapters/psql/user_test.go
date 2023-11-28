package psql_test

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/psql"
	"github.com/flowck/dobermann/backend/tests"
)

func TestUserRepository_Lifecycle(t *testing.T) {
	acc := tests.FixtureAndInsertAccount(t, db)
	user, err := acc.FirstAccountOwner()
	require.NoError(t, err)

	repo := psql.NewUserRepository(db)
	require.NoError(t, repo.Insert(ctx, user))

	userFound, err := repo.FindByEmail(ctx, user.Email())
	require.NoError(t, err)

	assert.Equal(t, user.ID(), userFound.ID())
}
