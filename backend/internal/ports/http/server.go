package http

import (
	"github.com/labstack/echo/v4"

	"github.com/flowck/doberman/internal/app"
)

type handlers struct {
	application *app.App
}

func RegisterHttpHandlers(application *app.App, router *echo.Group) {
	RegisterHandlers(router, &handlers{
		application: application,
	})
}
