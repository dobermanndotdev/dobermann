package command

import (
	"context"

	"github.com/flowck/doberman/internal/domain"
)

type ConfirmAccount struct {
	ConfirmationCode domain.ID
}

type confirmAccountUpdater interface {
	ConfirmAccount(ctx context.Context, confirmationCode domain.ID) error
}

type ConfirmAccountHandler struct {
	confirmAccountUpdater confirmAccountUpdater
}

func NewConfirmAccountHandler(confirmAccountUpdater confirmAccountUpdater) ConfirmAccountHandler {
	return ConfirmAccountHandler{
		confirmAccountUpdater: confirmAccountUpdater,
	}
}

func (c ConfirmAccountHandler) Execute(ctx context.Context, cmd ConfirmAccount) error {
	return c.confirmAccountUpdater.ConfirmAccount(ctx, cmd.ConfirmationCode)
}
