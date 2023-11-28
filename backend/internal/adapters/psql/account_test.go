package psql_test

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/adapters/psql"
	"github.com/flowck/dobermann/backend/tests"
)

func TestPsqlRepository_Insert(t *testing.T) {
	repo := psql.NewAccountRepository(db)
	require.NoError(t, repo.Insert(ctx, tests.FixtureAccount(t)))
}
