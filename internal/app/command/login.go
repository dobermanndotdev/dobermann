package command

import (
	"context"

	"github.com/dobermanndotdev/dobermann/internal/domain/account"
)

type LogIn struct {
	Email             account.Email
	PlainTextPassword string
}

type tokenSigner interface {
	Sign(metadata map[string]string) (string, error)
}

type LoginHandler struct {
	tokenSigner    tokenSigner
	userRepository account.UserRepository
}

func NewLoginHandler(userRepository account.UserRepository, tokenSigner tokenSigner) LoginHandler {
	return LoginHandler{
		tokenSigner:    tokenSigner,
		userRepository: userRepository,
	}
}

func (h LoginHandler) Execute(ctx context.Context, cmd LogIn) (string, error) {
	user, err := h.userRepository.FindByEmail(ctx, cmd.Email)
	if err != nil {
		return "", err
	}

	err = user.Authenticate(cmd.PlainTextPassword)
	if err != nil {
		return "", err
	}

	token, err := h.tokenSigner.Sign(map[string]string{
		"id":         user.ID().String(),
		"email":      user.Email().Address(),
		"account_id": user.AccountID().String(),
		"role":       user.Role().String(),
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
