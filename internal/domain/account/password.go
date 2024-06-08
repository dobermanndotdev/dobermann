package account

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/dobermanndotdev/dobermann/internal/common/hashing"
)

type Password struct {
	hash string
}

func NewPassword(plainText string) (Password, error) {
	plainText = strings.TrimSpace(plainText)

	if plainText == "" {
		return Password{}, errors.New("password cannot be empty")
	}

	hash, err := hashing.Hash(plainText)
	if err != nil {
		return Password{}, err
	}

	return Password{hash: hash}, nil
}

func NewPasswordFromHash(hash string) (Password, error) {
	hash = strings.TrimSpace(hash)
	_, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		return Password{}, errors.New("the hash provided is invalid")
	}

	return Password{hash: hash}, nil
}

func (p Password) IsEmpty() bool {
	return p.hash == ""
}

func (p Password) String() string {
	return p.hash
}
