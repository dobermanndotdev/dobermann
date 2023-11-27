package account

import "context"

type Repository interface {
	Insert(ctx context.Context, acc *Account) error
}

type UserRepository interface {
	FindByEmail(ctx context.Context, email Email) (*User, error)
	Insert(ctx context.Context, user *User) error
}
