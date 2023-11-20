package account

import (
	"context"

	"github.com/friendsofgo/errors"
)

var (
	ErrAccountUserExists = errors.New("user exists")
	ErrAccountNotFound   = errors.New("account not found")
)

type Repository interface {
	Insert(ctx context.Context, acc *Account) error
	Update(ctx context.Context, acc *Account) error
}
