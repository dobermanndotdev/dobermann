package app

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/app/command"
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
	CreateMonitor CommandHandler[command.CreateMonitor]
	CreateAccount CommandHandler[command.CreateAccount]
	LogIn         CommandHandlerWithResult[command.LogIn, string]
}

type App struct {
	Commands Commands
}
