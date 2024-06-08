package transaction

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/adapters/events"
	"github.com/dobermanndotdev/dobermann/internal/adapters/psql"
	"github.com/dobermanndotdev/dobermann/internal/app/command"
)

type PsqlProviderMock struct {
	Adapters *command.TransactableAdapters
}

func NewPsqlProviderMock() PsqlProviderMock {
	return PsqlProviderMock{
		Adapters: &command.TransactableAdapters{
			AccountRepository:  psql.NewAccountRepositoryMock(),
			UserRepository:     psql.NewUserRepositoryMock(),
			MonitorRepository:  psql.NewMonitorRepositoryMock(),
			IncidentRepository: psql.NewIncidentRepositoryMock(),
			EventPublisher:     events.NewPublisherMock(),
		},
	}
}

func (p PsqlProviderMock) Transact(ctx context.Context, f command.TransactFunc) error {
	if err := f(*p.Adapters); err != nil {
		return err
	}

	return nil
}
