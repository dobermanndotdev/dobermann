package httpport

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type HandlerError struct {
	*echo.HTTPError
	slug string
}

func NewError(err error, slug string) *HandlerError {
	return NewErrorWithStatus(err, slug, http.StatusInternalServerError)
}

func NewErrorWithStatus(err error, slug string, code int) *HandlerError {
	if err == nil {
		return nil
	}

	echoErr := echo.NewHTTPError(code, err.Error())
	echoErr = echoErr.SetInternal(err)

	return &HandlerError{
		HTTPError: echoErr,
		slug:      slug,
	}
}

func (e HandlerError) Slug() string {
	return e.slug
}

func (e HandlerError) Error() string {
	if e.Internal != nil {
		return e.Internal.Error()
	}

	return e.HTTPError.Error()
}
