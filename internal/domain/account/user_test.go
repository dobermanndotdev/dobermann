package account_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/account"
)

func mustEmail(t *testing.T) account.Email {
	email, err := account.NewEmail(gofakeit.Email())
	require.NoError(t, err)

	return email
}

func mustPassword(t *testing.T) account.Password {
	password, err := account.NewPassword(gofakeit.Password(
		true,
		true,
		true,
		true,
		false,
		16,
	))
	require.NoError(t, err)

	return password
}

func TestUser_NewUser(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr string
		id          domain.ID
		firstName   string
		lastName    string
		email       account.Email
		role        account.Role
		password    account.Password
		accountID   domain.ID
		createdAt   time.Time
	}{
		{
			name:        "new_user",
			expectedErr: "",
			id:          domain.NewID(),
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			email:       mustEmail(t),
			role:        account.RoleAdmin,
			password:    mustPassword(t),
			accountID:   domain.NewID(),
			createdAt:   time.Now(),
		},
		{
			name:        "error_invalid_id",
			expectedErr: "id cannot be invalid",
			id:          domain.ID{},
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			email:       mustEmail(t),
			role:        account.RoleAdmin,
			password:    mustPassword(t),
			accountID:   domain.NewID(),
			createdAt:   time.Now(),
		},
		{
			name:        "error_invalid_email_address",
			expectedErr: "email cannot be invalid",
			id:          domain.NewID(),
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			email:       account.Email{},
			role:        account.RoleAdmin,
			password:    mustPassword(t),
			accountID:   domain.NewID(),
			createdAt:   time.Now(),
		},
		{
			name:        "error_invalid_password",
			expectedErr: "password cannot be invalid",
			id:          domain.NewID(),
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			email:       mustEmail(t),
			role:        account.RoleAdmin,
			password:    account.Password{},
			accountID:   domain.NewID(),
			createdAt:   time.Now(),
		},
		{
			name:        "error_invalid_account_id",
			expectedErr: "",
			id:          domain.NewID(),
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			email:       mustEmail(t),
			role:        account.RoleAdmin,
			password:    mustPassword(t),
			accountID:   domain.ID{},
			createdAt:   time.Now(),
		},
		{
			name:        "error_created_at_set_in_the_future",
			expectedErr: "createdAt cannot be set in the future",
			id:          domain.NewID(),
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			email:       mustEmail(t),
			role:        account.RoleAdmin,
			password:    mustPassword(t),
			accountID:   domain.NewID(),
			createdAt:   time.Date(time.Now().Year()+1, 1, 1, 1, 1, 1, 0, time.UTC),
		},
		{
			name:        "error_role_is_invalid",
			expectedErr: "role cannot be invalid",
			id:          domain.NewID(),
			firstName:   gofakeit.FirstName(),
			lastName:    gofakeit.LastName(),
			email:       mustEmail(t),
			role:        account.Role{},
			password:    mustPassword(t),
			accountID:   domain.NewID(),
			createdAt:   time.Now(),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			user, err := account.NewUser(
				tc.id,
				tc.firstName,
				tc.lastName,
				tc.email,
				tc.role,
				tc.password,
				tc.accountID,
				tc.createdAt,
			)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}

			require.NotNil(t, user)
			require.NoError(t, err)

			assert.Equal(t, tc.id, user.ID())
			assert.Equal(t, tc.firstName, user.FirstName())
			assert.Equal(t, tc.lastName, user.LastName())
			assert.Equal(t, tc.email, user.Email())
			assert.Equal(t, tc.password, user.Password())
			assert.Equal(t, tc.accountID, user.AccountID())
			assert.Equal(t, tc.createdAt.UTC(), user.CreatedAt())
		})
	}
}
