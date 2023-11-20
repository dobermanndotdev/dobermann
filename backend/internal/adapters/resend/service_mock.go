package resend

import (
	"context"

	"github.com/flowck/doberman/internal/domain/account"
)

type ServiceMock struct{}

func (s ServiceMock) SendAccountConfirmationEmail(ctx context.Context, acc *account.Account) error {
	return nil
}
