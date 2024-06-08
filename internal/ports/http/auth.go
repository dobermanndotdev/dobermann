package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dobermanndotdev/dobermann/internal/app/command"
	"github.com/dobermanndotdev/dobermann/internal/domain/account"
)

func (h handlers) CreateAccount(c echo.Context) error {
	var body CreateAccountRequest
	if err := c.Bind(&body); err != nil {
		return NewHandlerError(err, "error-loading-the-payload")
	}

	email, err := account.NewEmail(body.Email)
	if err != nil {
		return NewHandlerErrorWithStatus(err, "validation-error", http.StatusBadRequest)
	}

	password, err := account.NewPassword(body.Password)
	if err != nil {
		return NewHandlerErrorWithStatus(err, "validation-error", http.StatusBadRequest)
	}

	newAccount, err := account.NewFirstTimeAccount(body.AccountName, email, password)
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

func (h handlers) Login(c echo.Context) error {
	var body LogInRequest
	if err := c.Bind(&body); err != nil {
		return NewHandlerError(err, "error-loading-the-payload")
	}

	email, err := account.NewEmail(body.Email)
	if err != nil {
		return NewHandlerErrorWithStatus(err, "validation-error", http.StatusBadRequest)
	}

	token, err := h.application.Commands.LogIn.Execute(c.Request().Context(), command.LogIn{
		Email:             email,
		PlainTextPassword: body.Password,
	})
	if errors.Is(err, account.ErrAuthenticationFailed) {
		return NewHandlerErrorWithStatus(err, "user-details-mismatch", http.StatusForbidden)
	}

	if err != nil {
		return NewHandlerError(err, "unable-to-authenticate")
	}

	return c.JSON(http.StatusOK, LogInPayload{Token: token})
}

func (h handlers) ConfirmInvitation(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}
