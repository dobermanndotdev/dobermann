package cqs

import "context"

type CommandHandler[C any] interface {
	Execute(ctx context.Context, cmd C) error
}

type CommandHandlerWithResult[C, R any] interface {
	Execute(ctx context.Context, cmd C) (R, err error)
}

type QueryHandler[Q, R any] interface {
	Execute(ctx context.Context, query Q) (R, err error)
}
