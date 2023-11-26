package account

import "errors"

var (
	ErrAccountExists = errors.New("we've found an account with the e-mail provided")
)
