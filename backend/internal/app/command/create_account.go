package command

import (
	"context"
	"fmt"

	"github.com/flowck/doberman/internal/domain/account"
)

type CreateAccount struct {
	Account *account.Account
}

type CreateAccountTxAdapters struct {
	AccountRepository account.Repository
	UserRepository    account.UserRepository
}

type CreateAccountHandler struct {
	txProvider TransactionProvider[CreateAccountTxAdapters]
}

func NewCreateAccountHandler(txProvider TransactionProvider[CreateAccountTxAdapters]) CreateAccountHandler {
	return CreateAccountHandler{
		txProvider: txProvider,
	}
}

func (h CreateAccountHandler) Execute(ctx context.Context, cmd CreateAccount) error {
	return h.txProvider.Transact(ctx, func(adapters CreateAccountTxAdapters) error {
		err := adapters.AccountRepository.Insert(ctx, cmd.Account)
		if err != nil {
			return fmt.Errorf("unable to save account: %v", err)
		}

		accountOwner, err := cmd.Account.FirstAccountOwner()
		if err != nil {
			return err
		}

		err = adapters.UserRepository.Insert(ctx, accountOwner)
		if err != nil {
			return fmt.Errorf("unable to save user: %v", err)
		}

		return nil
	})
}
