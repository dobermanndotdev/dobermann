package account

import "errors"

var (
	ErrAccountExists        = errors.New("we've found an account with the e-mail provided")
	ErrAuthenticationFailed = errors.New("the password provided doesn't match")
	ErrUserNotFound         = errors.New("user not found")
)
