package users

import (
	"github.com/volatiletech/null/v8"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain/account"
)

func mapUserToModel(user *account.User) models.User {
	return models.User{
		ID:        user.ID().String(),
		FirstName: null.StringFrom(user.FirstName()),
		LastName:  null.StringFrom(user.LastName()),
		Email:     user.Email(),
		Role:      user.Role().String(),
		AccountID: user.AccountID().String(),
		CreatedAt: user.CreatedAt(),
	}
}
