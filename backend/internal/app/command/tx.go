package command

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/domain/account"
)

type TransactableAdapters struct {
	AccountRepository account.Repository
	UserRepository    account.UserRepository
}

type TransactionProvider interface {
	Transact(ctx context.Context, f TransactFunc) error
}

type TransactFunc func(adapters TransactableAdapters) error
