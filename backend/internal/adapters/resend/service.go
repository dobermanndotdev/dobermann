package resend

import (
	"context"
	"fmt"
	"os"

	resendsdk "github.com/resendlabs/resend-go"

	"github.com/flowck/doberman/internal/common/ddd"
	"github.com/flowck/doberman/internal/domain/account"
)

type Service struct {
	client *resendsdk.Client
	from   ddd.Email
}

func NewService(client *resendsdk.Client, from ddd.Email) Service {
	return Service{
		client: client,
		from:   from,
	}
}

func (s Service) SendAccountConfirmationEmail(ctx context.Context, acc *account.Account) error {
	owner := acc.Users()[0]
	confirmationLink := fmt.Sprintf("%s?code=%s", os.Getenv("ACCOUNT_CONFIRMATION_LINK"), owner.ConfirmationCode())

	params := &resendsdk.SendEmailRequest{
		From:    s.from.Address(),
		To:      []string{acc.Users()[0].Email().Address()},
		Subject: "[Doberman] - Welcome",
		Html: fmt.Sprintf(`
Hi %s, </br>
Welcome to Doberman, please confirm your account by clicking <a href="%s">here</a>. </br></br>
Thank you
		`, owner.FirstName(), confirmationLink),
		Text: "",
	}

	_, err := s.client.Emails.Send(params)
	if err != nil {
		return fmt.Errorf("unable to send the account confirmation email: %v", err)
	}

	return nil
}
