package accounts

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/friendsofgo/errors"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/flowck/doberman/internal/adapters/models"
	"github.com/flowck/doberman/internal/domain"
	"github.com/flowck/doberman/internal/domain/account"
)

type PsqlRepository struct {
	db boil.ContextExecutor
}

func NewPsqlRepository(db boil.ContextExecutor) PsqlRepository {
	return PsqlRepository{db: db}
}

func (p PsqlRepository) Update(ctx context.Context, acc *account.Account) error {
	return nil
}

func (p PsqlRepository) Insert(ctx context.Context, acc *account.Account) error {
	model := mapAccountToModel(acc)

	err := model.Insert(ctx, p.db, boil.Infer())
	if err != nil {
		return err
	}

	owner := acc.Users()[0]
	err = model.AddUsers(ctx, p.db, true, &models.User{
		ID:               owner.ID().String(),
		AccountID:        acc.ID().String(),
		FirstName:        owner.FirstName(),
		LastName:         owner.LastName(),
		Email:            owner.Email().Address(),
		ConfirmationCode: null.StringFrom(owner.ConfirmationCode().String()),
	})
	if err != nil {
		return fmt.Errorf("unable to save owner: %v", err)
	}

	return nil
}

func (p PsqlRepository) ConfirmAccount(ctx context.Context, confirmationCode domain.ID) error {
	model, err := models.Users(models.UserWhere.ConfirmationCode.EQ(null.StringFrom(confirmationCode.String()))).One(ctx, p.db)
	if errors.Is(err, sql.ErrNoRows) {
		log.Println(err)
		return account.ErrAccountNotFound
	}

	if err != nil {
		return err
	}

	model.ConfirmationCode = null.StringFrom("")
	_, err = model.Update(ctx, p.db, boil.Infer())
	if err != nil {
		return err
	}

	return nil
}

func mapAccountToModel(acc *account.Account) *models.Account {
	return &models.Account{
		ID:   acc.ID().String(),
		Name: acc.Name(),
	}
}
