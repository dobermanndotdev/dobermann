package users

import (
	"context"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/doberman/internal/domain/account"
)

func NewPsqlRepository(db boil.ContextExecutor) PsqlRepository {
	return PsqlRepository{
		db: db,
	}
}

type PsqlRepository struct {
	db boil.ContextExecutor
}

func (p PsqlRepository) Insert(ctx context.Context, user *account.User) error {
	//TODO implement me
	panic("implement me")
}
