package http

import (
	"errors"
	"net/http"

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

	newMonitor, err := mapRequestToMonitor(body, user)
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

func (h handlers) GetMonitorByID(c echo.Context, monitorID string) error {
	_, err := retrieveUserFromCtx(c)
	if err != nil {
		return NewUnableToRetrieveUserFromCtx(err)
	}

	mID, err := domain.NewIdFromString(monitorID)
	if err != nil {
		return NewHandlerErrorWithStatus(err, "invalid-monitor-id", http.StatusBadRequest)
	}

	foundMonitor, err := h.application.Queries.MonitorByID.Execute(c.Request().Context(), query.MonitorByID{
		ID: mID,
	})
	if err != nil {
		return NewHandlerError(err, "unable-to-get-monitor")
	}

	return c.JSON(http.StatusOK, GetAllMonitorByIdPayload{
		Data: mapMonitorToResponseItem(foundMonitor),
	})
}

func (h handlers) ToggleMonitorPause(c echo.Context, monitorID string) error {
	_, err := retrieveUserFromCtx(c)
	if err != nil {
		return NewUnableToRetrieveUserFromCtx(err)
	}

	var body ToggleMonitorPauseRequest
	if err = c.Bind(&body); err != nil {
		return NewHandlerError(err, "error-loading-the-payload")
	}

	mID, err := domain.NewIdFromString(monitorID)
	if err != nil {
		return NewHandlerErrorWithStatus(err, "invalid-monitor-id", http.StatusBadRequest)
	}

	err = h.application.Commands.ToggleMonitorPause.Execute(c.Request().Context(), command.ToggleMonitorPause{
		MonitorID: mID,
		Pause:     body.Pause,
	})
	if err != nil {
		return NewHandlerError(err, "unable-to-get-monitor")
	}

	return c.NoContent(http.StatusNoContent)
}

func (h handlers) DeleteMonitor(c echo.Context, monitorID string) error {
	_, err := retrieveUserFromCtx(c)
	if err != nil {
		return NewUnableToRetrieveUserFromCtx(err)
	}

	mID, err := domain.NewIdFromString(monitorID)
	if err != nil {
		return NewHandlerErrorWithStatus(err, "invalid-monitor-id", http.StatusBadRequest)
	}

	err = h.application.Commands.DeleteMonitor.Execute(c.Request().Context(), command.DeleteMonitor{
		ID: mID,
	})
	if errors.Is(err, monitor.ErrMonitorNotFound) {
		return NewHandlerErrorWithStatus(err, "monitor-not-found", http.StatusNotFound)
	}

	if err != nil {
		return NewHandlerError(err, "unable-to-get-monitor")
	}

	return c.NoContent(http.StatusNoContent)
}
