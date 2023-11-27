package users

import (
	"github.com/volatiletech/null/v8"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/account"
)

func mapUserToModel(user *account.User) models.User {
	return models.User{
		ID:        user.ID().String(),
		FirstName: null.StringFrom(user.FirstName()),
		LastName:  null.StringFrom(user.LastName()),
		Email:     user.Email().Address(),
		Password:  user.Password().String(),
		Role:      user.Role().String(),
		AccountID: user.AccountID().String(),
		CreatedAt: user.CreatedAt(),
	}
}

func mapModelToUser(model *models.User) (*account.User, error) {
	id, err := domain.NewIdFromString(model.ID)
	if err != nil {
		return nil, err
	}

	accountId, err := domain.NewIdFromString(model.AccountID)
	if err != nil {
		return nil, err
	}

	email, err := account.NewEmail(model.Email)
	if err != nil {
		return nil, err
	}

	role, err := account.NewRole(model.Role)
	if err != nil {
		return nil, err
	}

	password, err := account.NewPasswordFromHash(model.Password)
	if err != nil {
		return nil, err
	}

	return account.NewUser(id, model.FirstName.String, model.LastName.String, email, role, password, accountId)
}
