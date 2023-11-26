package accounts

import (
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/flowck/dobermann/backend/internal/adapters/models"
	"github.com/flowck/dobermann/backend/internal/domain/account"
)

func mapAccountToModel(acc *account.Account) models.Account {
	return models.Account{
		ID:         acc.ID().String(),
		Name:       acc.Name(),
		VerifiedAt: null.TimeFromPtr(acc.VerifiedAt()),
		CreatedAt:  time.Now(),
	}
}
