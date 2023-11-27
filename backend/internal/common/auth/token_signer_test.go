package auth_test

import (
	"testing"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flowck/dobermann/backend/internal/common/auth"
)

func TestTokenSigner(t *testing.T) {
	signer, err := auth.NewTokenSigner(gofakeit.BeerName(), time.Hour)
	require.NoError(t, err)

	metadata := auth.Metadata{
		"id":    ulid.Make().String(),
		"email": gofakeit.Email(),
	}
	token, err := signer.Sign(metadata)
	require.NoError(t, err)

	verifiedMetadata, err := signer.Verify(token)
	require.NoError(t, err)

	assert.NotEmpty(t, token)
	assert.Equal(t, metadata["id"], verifiedMetadata["id"])
	assert.Equal(t, metadata["email"], verifiedMetadata["email"])

	_, err = signer.Verify(gofakeit.UUID())
	assert.Error(t, err)
}

func TestNewTokenSigner(t *testing.T) {
	_, err := auth.NewTokenSigner("", time.Hour)
	assert.EqualError(t, err, "the secret cannot be empty")

	_, err = auth.NewTokenSigner(gofakeit.BeerName(), time.Hour)
	require.NoError(t, err)
}
