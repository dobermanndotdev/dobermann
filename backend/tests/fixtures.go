package tests

import (
	"context"
	"database/sql"
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/require"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain/account"
)

func FixturePassword() string {
	return gofakeit.Password(true, true, true, true, false, 12)
}

func FixtureAccount(t *testing.T) *account.Account {
	email, err := account.NewEmail(gofakeit.Email())
	require.NoError(t, err)

	password, err := account.NewPassword(FixturePassword())
	require.NoError(t, err)

	acc, err := account.NewFirstTimeAccount(gofakeit.Company(), email, password)
	require.NoError(t, err)

	return acc
}

func FixtureAndInsertAccount(t *testing.T, db *sql.DB) *account.Account {
	acc := FixtureAccount(t)

	model := models.Account{
		Name: acc.Name(),
		ID:   acc.ID().String(),
	}

	require.NoError(t, model.Insert(context.Background(), db, boil.Infer()))

	return acc
}
