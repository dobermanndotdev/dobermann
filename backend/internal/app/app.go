package app

import (
	"context"

	"github.com/flowck/doberman/internal/app/command"
)

type QueryHandler[Q, R any] interface {
	Execute(ctx context.Context, q Q) (R, error)
}

type CommandHandler[C any] interface {
	Execute(ctx context.Context, cmd C) error
}

type CommandHandlerWithResult[C, R any] interface {
	Execute(ctx context.Context, cmd C) (result R, err error)
}

type Commands struct {
	CreateAccount CommandHandler[command.CreateAccount]
	CreateMonitor CommandHandler[any]
}

type App struct {
	Commands Commands
}
