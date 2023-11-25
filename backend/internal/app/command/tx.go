package command

import "context"

type TransactionProvider[T any] interface {
	Transact(ctx context.Context, f TransactFunc[T]) error
}

type TransactFunc[T any] func(adapters T) error
