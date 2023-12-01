package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/oklog/ulid/v2"
)

type TokenSigner struct {
	secret   []byte
	duration time.Duration
}

type Metadata = map[string]string

type customClaims struct {
	Metadata Metadata
	jwt.RegisteredClaims
}

func NewTokenSigner(secret string, duration time.Duration) (TokenSigner, error) {
	if secret == "" {
		return TokenSigner{}, errors.New("the secret cannot be empty")
	}

	return TokenSigner{
		duration: duration,
		secret:   []byte(secret),
	}, nil
}

func (s TokenSigner) Sign(metadata Metadata) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims{
		Metadata: metadata,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        ulid.Make().String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.duration)),
		},
	})

	return token.SignedString(s.secret)
}

func (s TokenSigner) Verify(token string) (Metadata, error) {
	t, err := jwt.ParseWithClaims(token, &customClaims{}, func(t *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := t.Claims.(*customClaims); ok && t.Valid {
		return claims.Metadata, nil
	}

	return nil, err
}
