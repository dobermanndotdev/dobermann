package http

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"github.com/flowck/doberman/internal/app/command"
)

func (h *handlers) CreateMonitor(c echo.Context) error {
	var body CreateMonitorRequest
	if err := c.Bind(&body); err != nil {
		return err
	}

	mo, err := mapReqBodyToMonitor(body)
	if err != nil {
		return err
	}

	err = h.application.CreateMonitor.Execute(c.Request().Context(), command.CreateMonitor{
		Monitor: mo,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusCreated)
}
