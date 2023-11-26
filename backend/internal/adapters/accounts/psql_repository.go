package accounts

import (
	"context"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/domain/account"
)

func NewPsqlRepository(db boil.ContextExecutor) PsqlRepository {
	return PsqlRepository{
		db: db,
	}
}

type PsqlRepository struct {
	db boil.ContextExecutor
}

func (p PsqlRepository) Insert(ctx context.Context, acc *account.Account) error {
	model := mapAccountToModel(acc)
	err := model.Insert(ctx, p.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("unable to save the account %s: %v", acc.ID(), err)
	}

	return nil
}
