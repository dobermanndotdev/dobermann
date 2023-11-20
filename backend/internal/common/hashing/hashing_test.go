package hashing_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/flowck/doberman/internal/common/hashing"
)

func TestHash(t *testing.T) {
	testCases := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "hash_password",
			password:    "my_super_secure_password",
			expectedErr: nil,
		},
		{
			name:        "attempt_to_hash_an_empty_string",
			password:    "",
			expectedErr: hashing.ErrPasswordTooShort,
		},
		{
			name:        "attempt_to_hash_a_short_password",
			password:    "hello_my",
			expectedErr: hashing.ErrPasswordTooShort,
		},
		{
			name:        "attempt_to_hash_a_password_with_sequential_empty_spaces",
			password:    "                     .       ",
			expectedErr: hashing.ErrPasswordTooShort,
		},
		{
			name:        "attempt_to_hash_a_long_password",
			password:    "8RD2djLQkWS6sUPTgAfC5WFUy2xS4ER9wBTFut83Huf99ce9cXbgBzQAhQxQghkuN7TEWZA72neQEcKXxTMuhB9v5pR7yVaZCr69bFRKJT9hKPjZgYaHkEMNxDaPn6Jyl",
			expectedErr: hashing.ErrPasswordTooLong,
		},
	}

	for i := range testCases {
		tCase := testCases[i]

		t.Run(tCase.name, func(t *testing.T) {
			_, err := hashing.Hash(tCase.password)
			assert.Equal(t, tCase.expectedErr, err)
		})
	}
}

func TestCompare(t *testing.T) {
	testCases := []struct {
		name           string
		password       string
		hashedPassword string
		expectedErr    error
	}{
		{
			name:           "compare_password",
			password:       "my_super_secure_password",
			hashedPassword: "$2a$11$Jc06agZKddnP1gcOTHAMuuQDrnP5b1.Nu2BrWqgV5XnbzRv7CB0d.",
			expectedErr:    nil,
		},
		{
			name:           "attempt_to_compare_an_empty_string",
			password:       "",
			hashedPassword: "$2a$11$Jc06agZKddnP1gcOTHAMuuQDrnP5b1.Nu2BrWqgV5XnbzRv7CB0d.",
			expectedErr:    hashing.ErrPasswordTooShort,
		},
		{
			name:           "attempt_to_compare_password",
			password:       "hello_my_password",
			hashedPassword: "$2a$11$Jc06agZKddnP1gcOTHAMuuQDrnP5b1.Nu2BrWqgV5XnbzRv7CB0d.",
			expectedErr:    hashing.ErrPasswordDontMatch,
		},
	}

	for i := range testCases {
		tCase := testCases[i]

		t.Run(tCase.name, func(t *testing.T) {
			err := hashing.Compare(tCase.password, tCase.hashedPassword)
			assert.Equal(t, tCase.expectedErr, err)
		})
	}
}
