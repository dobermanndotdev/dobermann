package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/flowck/dobermann/backend/internal/app/command"
	"github.com/flowck/dobermann/backend/internal/domain/account"
)

func (h handlers) CreateAccount(c echo.Context) error {
	var body CreateAccountRequest
	if err := c.Bind(&body); err != nil {
		return NewHandlerError(err, "error-loading-the-payload")
	}

	newAccount, err := account.NewFirstTimeAccount(body.AccountName, body.Email, body.Password)
	if err != nil {
		return NewHandlerErrorWithStatus(err, "validation-error", http.StatusBadRequest)
	}

	err = h.application.Commands.CreateAccount.Execute(c.Request().Context(), command.CreateAccount{
		Account: newAccount,
	})
	if errors.Is(err, account.ErrAccountExists) {
		return NewHandlerErrorWithStatus(err, "email-in-use", http.StatusBadRequest)
	}

	if err != nil {
		return NewHandlerError(err, "unable-to-create-account")
	}

	return c.NoContent(http.StatusCreated)
}

func (h handlers) Login(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}
