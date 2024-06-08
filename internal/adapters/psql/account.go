package psql

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/dobermanndotdev/dobermann/internal/domain/account"
)

func NewAccountRepository(db boil.ContextExecutor) AccountRepository {
	return AccountRepository{
		db: db,
	}
}

type AccountRepository struct {
	db boil.ContextExecutor
}

func (p AccountRepository) Insert(ctx context.Context, acc *account.Account) error {
	model := mapAccountToModel(acc)
	err := model.Insert(ctx, p.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("unable to save the account %s: %v", acc.ID(), err)
	}

	return nil
}
