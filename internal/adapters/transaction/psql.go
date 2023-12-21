package transaction

import (
	"context"
	"fmt"

	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/adapters/events"
	"github.com/flowck/dobermann/backend/internal/adapters/psql"
	"github.com/flowck/dobermann/backend/internal/app/command"
	"github.com/flowck/dobermann/backend/internal/common/logs"
)

type PsqlProvider struct {
	logger    *logs.Logger
	publisher message.Publisher
	db        boil.ContextBeginner
}

func NewPsqlProvider(db boil.ContextBeginner, publisher message.Publisher, logger *logs.Logger) PsqlProvider {
	return PsqlProvider{
		logger:    logger,
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
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			p.logger.WithError(err).WithField("rollback_err", rollbackErr).Error("Rollback error")
		}

		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
