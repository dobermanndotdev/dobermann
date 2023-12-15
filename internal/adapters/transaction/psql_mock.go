package transaction

import (
	"context"

	"github.com/flowck/dobermann/backend/internal/adapters/events"
	"github.com/flowck/dobermann/backend/internal/adapters/psql"
	"github.com/flowck/dobermann/backend/internal/app/command"
)

type PsqlProviderMock struct {
}

func NewPsqlProviderMock() PsqlProviderMock {
	return PsqlProviderMock{}
}

func (p PsqlProviderMock) Transact(ctx context.Context, f command.TransactFunc) error {
	adapters := command.TransactableAdapters{
		AccountRepository:  psql.NewAccountRepositoryMock(),
		UserRepository:     psql.NewUserRepositoryMock(),
		MonitorRepository:  psql.NewMonitorRepositoryMock(),
		IncidentRepository: psql.NewIncidentRepositoryMock(),
		EventPublisher:     events.NewPublisherMock(),
	}

	if err := f(adapters); err != nil {
		return err
	}

	return nil
}
