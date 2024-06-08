package account

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/domain"
)

type Repository interface {
	Insert(ctx context.Context, acc *Account) error
}

type UserRepository interface {
	Insert(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id domain.ID) (*User, error)
	FindByEmail(ctx context.Context, email Email) (*User, error)
}
