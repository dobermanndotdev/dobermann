package account

import "context"

type Repository interface {
	Insert(ctx context.Context, acc *Account) error
}

type UserRepository interface {
	Insert(ctx context.Context, user *User) error
}
