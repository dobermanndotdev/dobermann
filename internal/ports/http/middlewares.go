package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/dobermanndotdev/dobermann/internal/common/logs"
	"github.com/dobermanndotdev/dobermann/internal/common/observability"
)

func loggerMiddleware(logger *logs.Logger) echo.MiddlewareFunc {
	if logger.Level.String() == logs.DebugLevel.String() {
		return middleware.Logger()
	}

	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogMethod:    true,
		LogError:     false,
		LogRequestID: true,
		LogHeaders:   []string{echo.HeaderXCorrelationID},
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fields := logs.Fields{
				"method":         v.Method,
				"uri":            v.URI,
				"status":         v.Status,
				"correlation_id": v.RequestID,
			}

			if v.Error != nil || v.Status > http.StatusBadRequest {
				logger.WithFields(fields).Error("request handled with an error")
			} else {
				logger.WithFields(fields).Info("request handled successfully")
			}

			return nil
		},
	})
}

func correlationIdMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		correlationID := c.Request().Header.Get(echo.HeaderXCorrelationID)

		if correlationID == "" {
			correlationID = observability.NewCorrelationID()
		}

		ctx := observability.ContextWithCorrelationID(c.Request().Context(), correlationID)
		c.SetRequest(c.Request().WithContext(ctx))

		return next(c)
	}
}

func errorHandler(logger *logs.Logger) echo.HTTPErrorHandler {
	return func(err error, c echo.Context) {
		if c.Response().Committed {
			return
		}

		logger.WithError(err).Debug("Handler failed with error")

		code := http.StatusInternalServerError
		errorSlug := "internal-server-error"
		errorMessage := err.Error()

		switch e := err.(type) {
		case *echo.HTTPError:
			code = e.Code
			errorMessage = e.Error()
		case *HandlerError:
			code = e.Code
			errorSlug = e.Slug()
			errorMessage = e.Error()
		}

		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, ErrorResponse{
				Error:   errorSlug,
				Message: errorMessage,
			})
		}
		if err != nil {
			logger.WithError(err).Error("Failed to send error response")
		}
	}
}
