package http

import (
	"github.com/labstack/echo/v4"

	"github.com/flowck/dobermann/backend/internal/common/logs"
)

func (h handlers) CreateMonitor(c echo.Context) error {
	uID, err := retrieveUserIdFromCtx(c)
	if err != nil {
		return err
	}

	logs.Infof("User id: %s", uID)
	//TODO implement me
	panic("implement me")
}
