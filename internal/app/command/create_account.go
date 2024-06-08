package command

import (
	"context"
	"fmt"

	"github.com/dobermanndotdev/dobermann/internal/domain/account"
)

type CreateAccount struct {
	LoginProviderID string
	Email           string
}

type CreateAccountHandler struct {
	txProvider TransactionProvider
}

func NewCreateAccountHandler(txProvider TransactionProvider) CreateAccountHandler {
	return CreateAccountHandler{
		txProvider: txProvider,
	}
}

func (h CreateAccountHandler) Execute(ctx context.Context, cmd CreateAccount) error {
	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		acc := account.NewAccount()
		err := adapters.AccountRepository.Insert(ctx, acc)
		if err != nil {
			return fmt.Errorf("unable to save account: %v", err)
		}

		email, err := account.NewEmail(cmd.Email)
		if err != nil {
			return err
		}
		user, err := account.NewUser(account.RoleOwner, email, cmd.LoginProviderID, acc.ID())
		if err != nil {
			return err
		}

		// save user
		err = adapters.UserRepository.Insert(ctx, user)
		if err != nil {
			return fmt.Errorf("unable to save user: %v", err)
		}

		return nil
	})
}
