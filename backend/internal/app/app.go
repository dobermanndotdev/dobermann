package app

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/app/command"
	"github.com/flowck/dobermann/backend/internal/app/query"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
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
	CreateMonitor  CommandHandler[command.CreateMonitor]
	CreateAccount  CommandHandler[command.CreateAccount]
	CheckEndpoint  CommandHandler[command.CheckEndpoint]
	CreateIncident CommandHandler[command.CreateIncident]
	LogIn          CommandHandlerWithResult[command.LogIn, string]
}

type Queries struct {
	MonitorByID QueryHandler[query.MonitorByID, *monitor.Monitor]
	AllMonitors QueryHandler[query.AllMonitors, query.PaginatedResult[*monitor.Monitor]]
}

type App struct {
	Commands Commands
	Queries  Queries
}
