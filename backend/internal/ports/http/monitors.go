package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/flowck/dobermann/backend/internal/app/command"
	"github.com/flowck/dobermann/backend/internal/app/query"
	"github.com/flowck/dobermann/backend/internal/domain"
	"github.com/flowck/dobermann/backend/internal/domain/monitor"
)

func (h handlers) CreateMonitor(c echo.Context) error {
	user, err := retrieveUserFromCtx(c)
	if err != nil {
		return NewUnableToRetrieveUserFromCtx(err)
	}

	var body CreateMonitorRequest
	if err = c.Bind(&body); err != nil {
		return NewHandlerError(err, "error-loading-the-payload")
	}

	newMonitor, err := monitor.NewMonitor(
		domain.NewID(),
		body.EndpointUrl,
		user.AccountID,
		false,
		nil,
		time.Now().UTC(),
		nil,
	)
	if err != nil {
		return NewHandlerErrorWithStatus(err, "validation-error", http.StatusBadRequest)
	}

	err = h.application.Commands.CreateMonitor.Execute(c.Request().Context(), command.CreateMonitor{
		Monitor: newMonitor,
	})
	if err != nil {
		return NewHandlerError(err, "unable-to-create-monitor")
	}

	return c.NoContent(http.StatusCreated)
}

func (h handlers) GetAllMonitors(c echo.Context, params GetAllMonitorsParams) error {
	user, err := retrieveUserFromCtx(c)
	if err != nil {
		return NewUnableToRetrieveUserFromCtx(err)
	}

	result, err := h.application.Queries.AllMonitors.Execute(c.Request().Context(), query.AllMonitors{
		AccountID: user.AccountID,
		Params:    query.NewPaginationParams(params.Page, params.Limit),
	})
	if err != nil {
		return NewHandlerError(err, "unable-to-get-monitors")
	}

	return c.JSON(http.StatusOK, GetAllMonitorsPayload{
		Data:       mapMonitorsToResponseItems(result.Data),
		Page:       result.Page,
		PageCount:  result.PageCount,
		PerPage:    result.PerPage,
		TotalCount: result.TotalCount,
	})
}
