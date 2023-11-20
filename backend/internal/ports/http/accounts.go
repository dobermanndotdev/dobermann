package http

import (
	"net/http"

	"github.com/friendsofgo/errors"
	"github.com/labstack/echo/v4"

	"github.com/flowck/doberman/internal/app/command"
	"github.com/flowck/doberman/internal/common/ddd"
	"github.com/flowck/doberman/internal/common/httpport"
	"github.com/flowck/doberman/internal/domain"
	"github.com/flowck/doberman/internal/domain/account"
)

func (h *handlers) CreateAnAccount(c echo.Context) error {
	var body CreateAnAccountRequest
	if err := c.Bind(&body); err != nil {
		return httpport.NewError(err, "invalid-body")
	}

	acc, err := mapCreateAccountRequestToAccount(body)
	if err != nil {
		return httpport.NewError(err, "invalid-fields")
	}

	err = h.application.CreateAccount.Execute(c.Request().Context(), command.CreateAccount{
		Account: acc,
	})

	if errors.Is(err, account.ErrAccountUserExists) {
		return httpport.NewErrorWithStatus(err, "unable-to-create-an-account", http.StatusBadRequest)
	}

	if err != nil {
		return httpport.NewError(err, "unable-to-create-an-account")
	}

	return c.NoContent(http.StatusCreated)
}

func (h *handlers) ConfirmAccount(c echo.Context, confirmationCode string) error {
	confCode, err := domain.NewIdFromString(confirmationCode)
	if err != nil {
		return httpport.NewError(err, "invalid-fields")
	}

	err = h.application.ConfirmAccount.Execute(c.Request().Context(), command.ConfirmAccount{
		ConfirmationCode: confCode,
	})
	if errors.Is(err, account.ErrAccountNotFound) {
		return httpport.NewErrorWithStatus(err, "account-not-found", http.StatusNotFound)
	}

	if err != nil {
		return httpport.NewErrorWithStatus(err, "unable-to-confirm-account", http.StatusBadRequest)
	}

	return c.NoContent(http.StatusOK)
}

func mapCreateAccountRequestToAccount(body CreateAnAccountRequest) (*account.Account, error) {
	email, err := ddd.NewEmail(body.Email)
	if err != nil {
		return nil, err
	}

	password, err := account.NewPassword(body.Password)
	if err != nil {
		return nil, err
	}

	owner, err := account.NewUser(
		ddd.NewID(),
		body.FirstName,
		body.LastName,
		email,
		password,
		account.RoleOwner, domain.NewID())
	if err != nil {
		return nil, err
	}

	return account.NewAccountWithOwner(ddd.NewID(), body.AccountName, owner)
}
