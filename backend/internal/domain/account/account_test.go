package account_test

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/doberman/internal/common/ddd"
	"github.com/flowck/doberman/internal/domain"
	"github.com/flowck/doberman/internal/domain/account"
)

const fixtureWithOwner = true

func TestNewAccount(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr string

		id          ddd.ID
		accountName string
		users       []*account.User
	}{
		{
			name:        "create_an_account",
			expectedErr: "",

			id:          ddd.NewID(),
			accountName: gofakeit.Company(),
			users:       fixtureUsers(t, 10, fixtureWithOwner),
		},
		{
			name:        "error_invalid_id",
			expectedErr: "id cannot be invalid",

			id:          ddd.ID{},
			accountName: gofakeit.Company(),
			users:       fixtureUsers(t, 10, fixtureWithOwner),
		},
		{
			name:        "error_empty_account_name",
			expectedErr: "name cannot be empty",

			id:          ddd.NewID(),
			accountName: " ",
			users:       fixtureUsers(t, 10, fixtureWithOwner),
		},
		{
			name:        "error_no_users",
			expectedErr: "users cannot be empty",

			id:          ddd.NewID(),
			accountName: gofakeit.Company(),
			users:       []*account.User{},
		},
		{
			name:        "error_no_owner",
			expectedErr: "an account must have at least one user with the Owner role",

			id:          ddd.NewID(),
			accountName: gofakeit.Company(),
			users:       fixtureUsers(t, 10, false),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			newAccount, err := account.NewAccount(
				tc.id,
				tc.accountName,
				tc.users,
			)

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.id, newAccount.ID())
			assert.Equal(t, tc.accountName, newAccount.Name())
			assert.Equal(t, len(tc.users), len(newAccount.Users()))
		})
	}
}

func fixtureUsers(t *testing.T, userCount int, withOwner bool) []*account.User {
	roles := []account.Role{account.RoleWriter, account.RoleAdmin}
	users := make([]*account.User, userCount)

	if withOwner {
		roles = append(roles, account.RoleOwner)
	}

	for i := 0; i < userCount; i++ {
		user, err := account.NewUser(
			ddd.NewID(),
			gofakeit.FirstName(),
			gofakeit.LastName(),
			mustEmail(t, gofakeit.Email()),
			mustPassword(t),
			roles[gofakeit.Number(0, len(roles)-1)],
			domain.NewID(),
		)
		require.NoError(t, err)

		users[i] = user
	}

	return users
}
