package http

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/dobermanndotdev/dobermann/internal/app/query"
	"github.com/dobermanndotdev/dobermann/internal/domain/account"
)

func (h handlers) GetProfileDetails(c echo.Context) error {
	authUser, err := retrieveUserFromCtx(c)
	if err != nil {
		return NewUnableToRetrieveUserFromCtx(err)
	}

	ctx := c.Request().Context()
	user, err := h.application.Queries.UserByID.Execute(ctx, query.UserByID{
		ID: authUser.ID,
	})
	if errors.Is(err, account.ErrUserNotFound) {
		return NewHandlerErrorWithStatus(err, "user-not-found", http.StatusNotFound)
	}

	if err != nil {
		return NewHandlerError(err, "unable-to-get-profile-details")
	}

	return c.JSON(http.StatusOK, mapUserToResponse(user))
}

func (h handlers) BulkInviteMembersByEmail(ctx echo.Context) error {
	//TODO implement me
	panic("implement me")
}
