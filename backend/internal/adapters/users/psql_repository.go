package users

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
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

func (p PsqlRepository) FindByEmail(ctx context.Context, email account.Email) (*account.User, error) {
	model, err := models.Users(models.UserWhere.Email.EQ(email.Address())).One(ctx, p.db)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, account.ErrUserNotFound
	}

	if err != nil {
		return nil, fmt.Errorf("unable to find user by email %s: %v", email, err)
	}

	return mapModelToUser(model)
}

func (p PsqlRepository) Insert(ctx context.Context, user *account.User) error {
	model := mapUserToModel(user)
	err := model.Insert(ctx, p.db, boil.Infer())
	if err != nil {
		return fmt.Errorf("unable to save user %s: %v", model.ID, err)
	}

	return nil
}
