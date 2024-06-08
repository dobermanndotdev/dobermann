package app

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/app/command"
	"github.com/dobermanndotdev/dobermann/internal/app/query"
	"github.com/dobermanndotdev/dobermann/internal/domain/account"
	"github.com/dobermanndotdev/dobermann/internal/domain/monitor"
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
	// Monitor
	CreateMonitor                      CommandHandler[command.CreateMonitor]
	EditMonitor                        CommandHandler[command.EditMonitor]
	CheckEndpoint                      CommandHandler[command.CheckEndpoint]
	DeleteMonitor                      CommandHandler[command.DeleteMonitor]
	CreateIncident                     CommandHandler[command.CreateIncident]
	ResolveIncident                    CommandHandler[command.ResolveIncident]
	ToggleMonitorPause                 CommandHandler[command.ToggleMonitorPause]
	BulkCheckEndpoints                 CommandHandler[command.BulkCheckEndpoints]
	NotifyOnIncidentResolved           CommandHandler[command.NotifyOnIncidentResolved]
	NotifyMonitorSubscribersOnIncident CommandHandler[command.NotifyMonitorSubscribersOnIncident]

	// IAM
	CreateAccount CommandHandler[command.CreateAccount]
	LogIn         CommandHandlerWithResult[command.LogIn, string]
}

type Queries struct {
	MonitorByID              QueryHandler[query.MonitorByID, *monitor.Monitor]
	IncidentByID             QueryHandler[query.IncidentByID, *monitor.Incident]
	MonitorResponseTimeStats QueryHandler[query.MonitorResponseTimeStats, []query.ResponseTimeStat]
	AllMonitors              QueryHandler[query.AllMonitors, query.PaginatedResult[*monitor.Monitor]]
	AllIncidents             QueryHandler[query.AllIncidents, query.PaginatedResult[*monitor.Incident]]

	// IAM
	UserByID QueryHandler[query.UserByID, *account.User]
}

type App struct {
	Commands Commands
	Queries  Queries
}
