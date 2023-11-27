package transaction

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/adapters/accounts"
	"github.com/flowck/dobermann/backend/internal/adapters/events"
	"github.com/flowck/dobermann/backend/internal/adapters/monitors"
	"github.com/flowck/dobermann/backend/internal/adapters/users"
	"github.com/flowck/dobermann/backend/internal/app/command"
)

type PsqlProvider struct {
	db        boil.ContextBeginner
	publisher message.Publisher
}

func NewPsqlProvider(db boil.ContextBeginner, publisher message.Publisher) PsqlProvider {
	return PsqlProvider{
		db:        db,
		publisher: publisher,
	}
}

func (p PsqlProvider) Transact(ctx context.Context, f command.TransactFuncc) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %v", err)
	}

	adapters := command.TransactableAdapters{
		AccountRepository: accounts.NewPsqlRepository(tx),
		UserRepository:    users.NewPsqlRepository(tx),
		MonitorRepository: monitors.NewPsqlRepository(tx),
		EventPublisher:    events.NewPublisher(p.publisher),
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
