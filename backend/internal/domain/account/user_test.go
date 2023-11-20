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

func TestNewUser(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr string

		id               ddd.ID
		firstName        string
		lastName         string
		email            ddd.Email
		password         account.Password
		role             account.Role
		confirmationCode domain.ID
	}{
		{
			name:        "create_a_user",
			expectedErr: "",

			id:               ddd.NewID(),
			firstName:        gofakeit.FirstName(),
			lastName:         gofakeit.LastName(),
			email:            mustEmail(t, gofakeit.Email()),
			password:         mustPassword(t),
			role:             account.RoleAdmin,
			confirmationCode: domain.NewID(),
		},
		{
			name:        "error_invalid_id",
			expectedErr: "id cannot be invalid",

			id:               ddd.ID{},
			firstName:        gofakeit.FirstName(),
			lastName:         gofakeit.LastName(),
			email:            mustEmail(t, gofakeit.Email()),
			password:         mustPassword(t),
			role:             account.RoleAdmin,
			confirmationCode: domain.NewID(),
		},
		{
			name:        "error_empty_first_name",
			expectedErr: "firstName cannot be empty",

			id:               ddd.NewID(),
			firstName:        " ",
			lastName:         gofakeit.LastName(),
			email:            mustEmail(t, gofakeit.Email()),
			password:         mustPassword(t),
			role:             account.RoleAdmin,
			confirmationCode: domain.NewID(),
		},
		{
			name:        "error_empty_last_name",
			expectedErr: "lastName cannot be empty",

			id:               ddd.NewID(),
			firstName:        gofakeit.FirstName(),
			lastName:         "     ",
			email:            mustEmail(t, gofakeit.Email()),
			password:         mustPassword(t),
			role:             account.RoleAdmin,
			confirmationCode: domain.NewID(),
		},
		{
			name:        "error_invalid_email",
			expectedErr: "",

			id:               ddd.NewID(),
			firstName:        gofakeit.FirstName(),
			lastName:         gofakeit.LastName(),
			email:            ddd.Email{},
			password:         mustPassword(t),
			role:             account.RoleAdmin,
			confirmationCode: domain.NewID(),
		},
		{
			name:        "error_invalid_password",
			expectedErr: "password cannot be invalid",

			id:               ddd.NewID(),
			firstName:        gofakeit.FirstName(),
			lastName:         gofakeit.LastName(),
			email:            mustEmail(t, gofakeit.Email()),
			password:         account.Password{},
			role:             account.RoleAdmin,
			confirmationCode: domain.NewID(),
		},
		{
			name:        "error_invalid_password",
			expectedErr: "role cannot be invalid",

			id:               ddd.NewID(),
			firstName:        gofakeit.FirstName(),
			lastName:         gofakeit.LastName(),
			email:            mustEmail(t, gofakeit.Email()),
			password:         mustPassword(t),
			role:             account.Role{},
			confirmationCode: domain.NewID(),
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
				tc.password,
				tc.role,
				tc.confirmationCode,
			)

			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.id, user.ID())
			assert.Equal(t, tc.firstName, user.FirstName())
			assert.Equal(t, tc.lastName, user.LastName())
			assert.Equal(t, tc.email, user.Email())
			assert.Equal(t, tc.password, user.Password())
			assert.Equal(t, tc.confirmationCode, user.ConfirmationCode())
		})
	}
}

func mustPassword(t *testing.T) account.Password {
	password, err := account.NewPassword(
		gofakeit.Password(true, true, true, true, false, 12),
	)
	require.NoError(t, err)

	return password
}

func mustEmail(t *testing.T, address string) ddd.Email {
	email, err := ddd.NewEmail(address)
	require.NoError(t, err)

	return email
}
