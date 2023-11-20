package app

import (
	"context"

	"github.com/flowck/doberman/internal/app/command"
)

type CommandHandler[C any] interface {
	Execute(ctx context.Context, cmd C) error
}

type QueryHandler[Q any, R any] interface {
	Execute(ctx context.Context, q Q) error
}

type Commands struct {
	CreateMonitor                          CommandHandler[command.CreateMonitor]
	EnqueueMonitors                        CommandHandler[command.EnqueueMonitors]
	CheckMonitorEndpoint                   CommandHandler[command.CheckMonitorEndpoint]
	ResolveIncident                        CommandHandler[any]
	CreateIncident                         CommandHandler[any]
	AcknowledgeIncident                    CommandHandler[any]
	EscalateIncident                       CommandHandler[any]
	NotifyTeamMembersOnCall                CommandHandler[any]
	NotifyTeamMembersInTheEscalationPolicy CommandHandler[any]
	PostCommentOnIncident                  CommandHandler[any]

	// Accounts
	CreateAccount  CommandHandler[command.CreateAccount]
	ConfirmAccount CommandHandler[command.ConfirmAccount]
}

type Queries struct{}

type App struct {
	Queries
	Commands
}
