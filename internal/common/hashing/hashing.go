package hashing

import (
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost     int8  = 11
	passwordMinLen int8  = 12
	passwordMaxLen int16 = 128
)

var (
	ErrPasswordTooShort  = fmt.Errorf("password cannot be less than %d", passwordMinLen)
	ErrPasswordTooLong   = fmt.Errorf("password cannot be greater than %d", passwordMaxLen)
	ErrPasswordDontMatch = errors.New("password provided doesn't match")
)

func Hash(password string) (string, error) {
	password = strings.TrimSpace(password)

	if err := checkPasswordLen(password); err != nil {
		return "", err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), int(bcryptCost))
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func Compare(password, hashedPassword string) error {
	if err := checkPasswordLen(password); err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return ErrPasswordDontMatch
	}

	return nil
}

func IsHash(hash string) bool {
	if _, err := bcrypt.Cost([]byte(hash)); err != nil {
		return false
	}

	return true
}

func checkPasswordLen(password string) error {
	if len(password) < int(passwordMinLen) {
		return ErrPasswordTooShort
	}

	if len(password) > int(passwordMaxLen) {
		return ErrPasswordTooLong
	}

	return nil
}
