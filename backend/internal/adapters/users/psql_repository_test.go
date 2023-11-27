package users_test

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/users"
	"github.com/flowck/dobermann/backend/internal/common/psql"
	"github.com/flowck/dobermann/backend/tests"
)

var (
	db  *sql.DB
	ctx context.Context
)

func TestPsqlRepository_Lifecycle(t *testing.T) {
	acc := tests.FixtureAndInsertAccount(t, db)
	user, err := acc.FirstAccountOwner()
	require.NoError(t, err)

	repo := users.NewPsqlRepository(db)
	require.NoError(t, repo.Insert(ctx, user))

	userFound, err := repo.FindByEmail(ctx, user.Email())
	require.NoError(t, err)

	assert.Equal(t, user.ID(), userFound.ID())
}

func TestMain(m *testing.M) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(context.Background(), time.Minute*1)
	defer cancel()

	var err error
	db, err = psql.Connect(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	os.Exit(m.Run())
}
