package psql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/dobermanndotdev/dobermann/internal/adapters/models"
	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/account"
)

func NewUserRepository(db boil.ContextExecutor) UserRepository {
	return UserRepository{
		db: db,
	}
}

type UserRepository struct {
	db boil.ContextExecutor
}

func (p UserRepository) FindByID(ctx context.Context, id domain.ID) (*account.User, error) {
	model, err := models.FindUser(ctx, p.db, id.String())
	if errors.Is(err, sql.ErrNoRows) {
		return nil, account.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return mapModelToUser(model)
}

func (p UserRepository) FindByEmail(ctx context.Context, email account.Email) (*account.User, error) {
	model, err := models.Users(models.UserWhere.Email.EQ(email.Address())).One(ctx, p.db)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, account.ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return mapModelToUser(model)
}

func (p UserRepository) Insert(ctx context.Context, user *account.User) error {
	exists, err := models.Users(models.UserWhere.Email.EQ(user.Email().Address())).Exists(ctx, p.db)
	if err != nil {
		return fmt.Errorf("unable to check if %s is already taken: %v", user.Email(), err)
	}

	if exists {
		return account.ErrAccountExists
	}

	model := mapUserToModel(user)
	return model.Insert(ctx, p.db, boil.Infer())
}
