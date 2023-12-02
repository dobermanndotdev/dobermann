package transaction

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/adapters/events"
	"github.com/flowck/dobermann/backend/internal/adapters/psql"
	"github.com/flowck/dobermann/backend/internal/app/command"
)

type PsqlProvider struct {
	publisher message.Publisher
	db        boil.ContextBeginner
}

func NewPsqlProvider(db boil.ContextBeginner, publisher message.Publisher) PsqlProvider {
	return PsqlProvider{
		db:        db,
		publisher: publisher,
	}
}

func (p PsqlProvider) Transact(ctx context.Context, f command.TransactFunc) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %v", err)
	}

	adapters := command.TransactableAdapters{
		AccountRepository:  psql.NewAccountRepository(tx),
		UserRepository:     psql.NewUserRepository(tx),
		MonitorRepository:  psql.NewMonitorRepository(tx),
		IncidentRepository: psql.NewIncidentRepository(tx),
		EventPublisher:     events.NewPublisher(p.publisher),
	}

	if err = f(adapters); err != nil {
		if rollbackErr := tx.Rollback(); err != nil {
			return fmt.Errorf("an error while trying to rollback transaction failed due to: %v %v", err, rollbackErr)
		}

		return fmt.Errorf("transaction rolled back due to: %v", err)
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("unable to commit transaction: %v", err)
	}

	return nil
}
