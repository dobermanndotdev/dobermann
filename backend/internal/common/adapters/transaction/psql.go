package transaction

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"

	commontx "github.com/flowck/doberman/internal/common/tx"
)

type PsqlProvider[T any] struct {
	db       boil.ContextBeginner
	adapters T
}

func NewPsqlProvider[T any](db boil.ContextBeginner, adapters T) PsqlProvider[T] {
	return PsqlProvider[T]{
		db:       db,
		adapters: adapters,
	}
}

func (p PsqlProvider[T]) Transact(ctx context.Context, f commontx.TransactFunc[T]) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("unable to begin transaction: %v", err)
	}

	if err = f(p.adapters); err != nil {
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
