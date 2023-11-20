package command

import (
	"context"

	"github.com/flowck/doberman/internal/common/tx"
	"github.com/flowck/doberman/internal/domain/account"
)

type CreateAccount struct {
	Account *account.Account
}

type notificationSender interface {
	SendAccountConfirmationEmail(ctx context.Context, acc *account.Account) error
}

type CreateAccountHandler struct {
	txProvider         tx.TransactionProvider[TransactableAdapters]
	notificationSender notificationSender
}

func NewCreateAccountHandler(
	txProvider tx.TransactionProvider[TransactableAdapters],
	notificationSender notificationSender,
) CreateAccountHandler {
	return CreateAccountHandler{
		txProvider:         txProvider,
		notificationSender: notificationSender,
	}
}

func (h CreateAccountHandler) Execute(ctx context.Context, cmd CreateAccount) error {
	return h.txProvider.Transact(ctx, func(adapters TransactableAdapters) error {
		err := adapters.AccountRepository.Insert(ctx, cmd.Account)
		if err != nil {
			return err
		}

		err = adapters.EventPublisher.PublishAccountCreated(ctx, AccountCreatedEvent{
			ID: cmd.Account.ID().String(),
		})
		if err != nil {
			return err
		}

		err = h.notificationSender.SendAccountConfirmationEmail(ctx, cmd.Account)
		if err != nil {
			return err
		}

		return nil
	})
}
