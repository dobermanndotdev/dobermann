package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"

	"github.com/dobermanndotdev/dobermann/internal/domain"
	"github.com/dobermanndotdev/dobermann/internal/domain/account"
)

const ctxAuthenticatedUser = "authenticated_user"

var (
	errMissingAuthorizationToken = errors.New("missing authorization token")
)

type jwtVerifier interface {
	Verify(token string) (map[string]string, error)
}

type authenticatedUser struct {
	ID        domain.ID
	AccountID domain.ID
	Email     account.Email
	Role      account.Role
}

func NewAuthenticator(verifier jwtVerifier) openapi3filter.AuthenticationFunc {
	return func(ctx context.Context, input *openapi3filter.AuthenticationInput) error {
		return Authenticate(verifier, ctx, input)
	}
}

func Authenticate(verifier jwtVerifier, ctx context.Context, input *openapi3filter.AuthenticationInput) error {
	authorizationHeader := input.RequestValidationInput.Request.Header.Get(echo.HeaderAuthorization)
	if authorizationHeader == "" {
		return NewHandlerErrorWithStatus(errMissingAuthorizationToken, "missing-authorization-token", http.StatusUnauthorized)
	}

	metadata, err := verifier.Verify(strings.TrimPrefix(authorizationHeader, "Bearer "))
	if err != nil {
		return NewHandlerErrorWithStatus(err, "invalid-authorization-token", http.StatusUnauthorized)
	}

	eCtx := oapimiddleware.GetEchoContext(ctx)
	eCtx.Set(ctxAuthenticatedUser, metadata)

	return nil
}

func retrieveUserFromCtx(c echo.Context) (*authenticatedUser, error) {
	metadata := c.Get(ctxAuthenticatedUser).(map[string]string)
	if metadata == nil {
		return nil, errors.New("unable to retrieve user from request context")
	}

	userID, err := domain.NewIdFromString(metadata["id"])
	if err != nil {
		return nil, fmt.Errorf("user retrieval from request context:%v", err)
	}

	accountID, err := domain.NewIdFromString(metadata["account_id"])
	if err != nil {
		return nil, fmt.Errorf("user retrieval from request context:%v", err)
	}

	email, err := account.NewEmail(metadata["email"])
	if err != nil {
		return nil, fmt.Errorf("user retrieval from request context:%v", err)
	}

	role, err := account.NewRole(metadata["role"])
	if err != nil {
		return nil, fmt.Errorf("user retrieval from request context:%v", err)
	}

	return &authenticatedUser{
		ID:        userID,
		Email:     email,
		AccountID: accountID,
		Role:      role,
	}, nil
}

/*func retrieveUserIdFromCtx(c echo.Context) (domain.ID, error) {
	data, err := retrieveUserFromCtx(c)
	if err != nil {
		return domain.ID{}, err
	}

	return data.ID, nil
}*/
