package psql_test

import (
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"github.com/dobermanndotdev/dobermann/internal/adapters/psql"
	"github.com/dobermanndotdev/dobermann/tests"
)

func TestPsqlRepository_Insert(t *testing.T) {
	repo := psql.NewAccountRepository(db)
	require.NoError(t, repo.Insert(ctx, tests.FixtureAccount(t)))
}
