package account_test

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dobermanndotdev/dobermann/internal/domain/account"
)

func TestNewEmail(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr string
		address     string
	}{
		{
			name:        "create_new_email",
			expectedErr: "",
			address:     gofakeit.Email(),
		},
		{
			name:        "error_empty_email_address",
			expectedErr: "address cannot be empty",
			address:     "",
		},
		{
			name:        "error_leading_and_trailing_spaces",
			expectedErr: "address cannot be empty",
			address:     "      ",
		},
		{
			name:        "error_invalid_email_missing_host",
			expectedErr: "the address provided is invalid",
			address:     "hello@",
		},
		{
			name:        "error_invalid_email_missing_username",
			expectedErr: "the address provided is invalid",
			address:     "@gmail.com",
		},
		{
			name:        "error_invalid_email_missing_at_character",
			expectedErr: "the address provided is invalid",
			address:     "hellogmail.com",
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			email, err := account.NewEmail(tc.address)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.address, email.Address())
		})
	}
}
