package domain_test

import (
	"testing"

	"github.com/brianvoe/gofakeit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/domain"
)

func TestNewIdFromString(t *testing.T) {
	testCases := []struct {
		name        string
		expectedErr string
		id          string
	}{
		{
			name:        "create_new_id_from_string",
			expectedErr: "",
			id:          domain.NewID().String(),
		},
		{
			name:        "error_empty_id",
			expectedErr: "id cannot be empty",
			id:          "       ",
		},
		{
			name:        "error_invalid_id_string",
			expectedErr: "id cannot be invalid",
			id:          gofakeit.Sentence(4),
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			id, err := domain.NewIdFromString(tc.id)
			if tc.expectedErr != "" {
				assert.EqualError(t, err, tc.expectedErr)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tc.id, id.String())
		})
	}
}
