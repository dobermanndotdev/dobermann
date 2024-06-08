package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dobermanndotdev/dobermann/internal/app/command"
)

func (h handlers) CreateAccount(c echo.Context) error {
	var body CreateAccountRequest
	if err := c.Bind(&body); err != nil {
		return NewHandlerError(err, "error-loading-the-payload")
	}

	if body.Data.EmailAddresses == nil || len(*body.Data.EmailAddresses) == 0 {
		return NewHandlerErrorWithStatus(errors.New("missing email"), "missing-email", http.StatusBadRequest)
	}

	emailAddresses := *body.Data.EmailAddresses
	email := emailAddresses[0].EmailAddress

	err := h.application.Commands.CreateAccount.Execute(c.Request().Context(), command.CreateAccount{
		Email:           *email,
		LoginProviderID: body.Data.Id,
	})
	if err != nil {
		return NewHandlerError(err, "unable-to-create-account")
	}

	return c.NoContent(http.StatusCreated)
}

func (h handlers) ConfirmInvitation(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}
