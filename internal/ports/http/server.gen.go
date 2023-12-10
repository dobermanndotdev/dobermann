// Package http provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package http

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/labstack/echo/v4"
)

const (
	BearerAuthScopes = "BearerAuth.Scopes"
)

// CreateAccountRequest defines model for CreateAccountRequest.
type CreateAccountRequest struct {
	AccountName string `json:"account_name"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

// CreateMonitorRequest defines model for CreateMonitorRequest.
type CreateMonitorRequest struct {
	EndpointUrl string `json:"endpoint_url"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	// Error Error custom error code such as 'email_in_use'
	Error string `json:"error"`

	// Message A description about the error
	Message string `json:"message"`
}

// GetAllMonitorByIdPayload defines model for GetAllMonitorByIdPayload.
type GetAllMonitorByIdPayload struct {
	Data Monitor `json:"data"`
}

// GetAllMonitorsPayload defines model for GetAllMonitorsPayload.
type GetAllMonitorsPayload struct {
	Data       []Monitor `json:"data"`
	Page       int       `json:"page"`
	PageCount  int       `json:"page_count"`
	PerPage    int       `json:"per_page"`
	TotalCount int64     `json:"total_count"`
}

// Incident defines model for Incident.
type Incident struct {
	CreatedAt time.Time `json:"created_at"`
	Id        string    `json:"id"`
}

// LogInPayload defines model for LogInPayload.
type LogInPayload struct {
	Token string `json:"token"`
}

// LogInRequest defines model for LogInRequest.
type LogInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Monitor defines model for Monitor.
type Monitor struct {
	CreatedAt     time.Time  `json:"created_at"`
	EndpointUrl   string     `json:"endpoint_url"`
	Id            string     `json:"id"`
	Incidents     []Incident `json:"incidents"`
	IsEndpointUp  bool       `json:"is_endpoint_up"`
	IsPaused      bool       `json:"is_paused"`
	LastCheckedAt *time.Time `json:"last_checked_at,omitempty"`
}

// ToggleMonitorPauseRequest defines model for ToggleMonitorPauseRequest.
type ToggleMonitorPauseRequest struct {
	Pause bool `json:"pause"`
}

// DefaultError defines model for DefaultError.
type DefaultError = ErrorResponse

// GetAllMonitorsParams defines parameters for GetAllMonitors.
type GetAllMonitorsParams struct {
	Page  *int `form:"page,omitempty" json:"page,omitempty"`
	Limit *int `form:"limit,omitempty" json:"limit,omitempty"`
}

// CreateAccountJSONRequestBody defines body for CreateAccount for application/json ContentType.
type CreateAccountJSONRequestBody = CreateAccountRequest

// LoginJSONRequestBody defines body for Login for application/json ContentType.
type LoginJSONRequestBody = LogInRequest

// CreateMonitorJSONRequestBody defines body for CreateMonitor for application/json ContentType.
type CreateMonitorJSONRequestBody = CreateMonitorRequest

// ToggleMonitorPauseJSONRequestBody defines body for ToggleMonitorPause for application/json ContentType.
type ToggleMonitorPauseJSONRequestBody = ToggleMonitorPauseRequest

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Creates a new account
	// (POST /auth/accounts)
	CreateAccount(ctx echo.Context) error
	// Log in
	// (POST /auth/login)
	Login(ctx echo.Context) error
	// Get all monitors in a with pagination
	// (GET /monitors)
	GetAllMonitors(ctx echo.Context, params GetAllMonitorsParams) error
	// Create a new monitor
	// (POST /monitors)
	CreateMonitor(ctx echo.Context) error
	// Get all monitors in a with pagination
	// (GET /monitors/{monitorID})
	GetMonitorByID(ctx echo.Context, monitorID string) error
	// Pause or unpause the monitor
	// (POST /monitors/{monitorID})
	ToggleMonitorPause(ctx echo.Context, monitorID string) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// CreateAccount converts echo context to params.
func (w *ServerInterfaceWrapper) CreateAccount(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateAccount(ctx)
	return err
}

// Login converts echo context to params.
func (w *ServerInterfaceWrapper) Login(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.Login(ctx)
	return err
}

// GetAllMonitors converts echo context to params.
func (w *ServerInterfaceWrapper) GetAllMonitors(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Parameter object where we will unmarshal all parameters from the context
	var params GetAllMonitorsParams
	// ------------- Optional query parameter "page" -------------

	err = runtime.BindQueryParameter("form", true, false, "page", ctx.QueryParams(), &params.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter page: %s", err))
	}

	// ------------- Optional query parameter "limit" -------------

	err = runtime.BindQueryParameter("form", true, false, "limit", ctx.QueryParams(), &params.Limit)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter limit: %s", err))
	}

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetAllMonitors(ctx, params)
	return err
}

// CreateMonitor converts echo context to params.
func (w *ServerInterfaceWrapper) CreateMonitor(ctx echo.Context) error {
	var err error

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.CreateMonitor(ctx)
	return err
}

// GetMonitorByID converts echo context to params.
func (w *ServerInterfaceWrapper) GetMonitorByID(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "monitorID" -------------
	var monitorID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "monitorID", runtime.ParamLocationPath, ctx.Param("monitorID"), &monitorID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter monitorID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.GetMonitorByID(ctx, monitorID)
	return err
}

// ToggleMonitorPause converts echo context to params.
func (w *ServerInterfaceWrapper) ToggleMonitorPause(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "monitorID" -------------
	var monitorID string

	err = runtime.BindStyledParameterWithLocation("simple", false, "monitorID", runtime.ParamLocationPath, ctx.Param("monitorID"), &monitorID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter monitorID: %s", err))
	}

	ctx.Set(BearerAuthScopes, []string{""})

	// Invoke the callback with all the unmarshalled arguments
	err = w.Handler.ToggleMonitorPause(ctx, monitorID)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.POST(baseURL+"/auth/accounts", wrapper.CreateAccount)
	router.POST(baseURL+"/auth/login", wrapper.Login)
	router.GET(baseURL+"/monitors", wrapper.GetAllMonitors)
	router.POST(baseURL+"/monitors", wrapper.CreateMonitor)
	router.GET(baseURL+"/monitors/:monitorID", wrapper.GetMonitorByID)
	router.POST(baseURL+"/monitors/:monitorID", wrapper.ToggleMonitorPause)

}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/9xY3W7jNhN9FYLfB+RGtZzuolj4qk7sBm6T3SBJ0YvAMGhqLHEjkgo5StYI/O4FKcmS",
	"LDnxos4W6J0ZkmdmzvzwKC+Ua5lpBQotHb1QAzbTyoJfTGDF8hSnxmjj1lwrBIXuJ8uyVHCGQqvwq9XK",
	"/c3yBCRzv/5vYEVH9H9hDR4Wuzb0aDelGbrZbAIageVGZA6MjuiYxKDACE7AHSWmPhuUNrx35wYYwphz",
	"nSu8gcccrHctMzoDg6KIgRX7C8UkuDWuM6AjatEIFdNNQEEykfbuZMzaZ22ins1NQA085sJAREf3bSMV",
	"ZANgvglKd6+0Euji3+MuqCjTQuEiN+nbdlunnZE2uV30KpNtxv0twnOLWpakcx0BsTlPCLPkxAe0EGqR",
	"WzihQZcqCdayGLrQY9JYE7bUORJMoLDSRdqNrzxVwc+3F/TyK3B0pi8Ax2la8nq2nkXXbJ1qFnWjjxi+",
	"WZ0lTscTf3e+a86+aUsgSHuw0W14zBi2LmowbpatUAgxmGpn4ctuzz6Yxf7bqJGl9fWVNpJhceSXj3Ve",
	"tjd26PDADRstd9roQc3dTHERlQOkTRf3zREtWNubiCH8hMI3VafmxAGNKSIaNMGdF5c6nqm9iUP9AOpt",
	"4OLYFm5/Px9huvSOk6pmjkLlG1NnD9cBFWVC7cGlvi2BnloXdlE7kjUMLrVOganyTMZyC1H/dsosLngC",
	"/OG7GOirmhYnHeeanjR56FTbnY7jtBr71+7C3lrxcH1xdZrPnZv75xB4bgSubx29BcoZMANmnGPiVku/",
	"+q3i4Pe/7mj5iHoLfrfmI0HMiicZviEYxdKJ5rY71d05OwrDWGCSLwdcy3CV6mf+EEZ6CUYypcKb6Xhy",
	"NR1Ix48vq0NuFUW10pXaYBwbXURXwkih9IAnTMVMiV9jt+GQaEdGTCrME0uWjD+Acp6kgkP5NBaKgF7N",
	"7r7Hw/Bydj79fOsDczUMRtovq1swT4LDgUEGFAWm7nQNW7v4BMYWIQwHw8Gps6IzUCwTdEQ/DIaDD34W",
	"YOITE7Ick7AUIEUZ6aK4XGl5gTaL6KitlmhRUGDxTEfro0m7XkW2aZcvmhz8Hxo68+fhabfGzm+m47vp",
	"pEis16H7zG+xwpZg9f2RS8nMehu/JYwoeCZsywOy2Lqm8g3je6pgNNWxUPvpvPTb70Nj60U5iL7hcW1X",
	"j2OPOL+dXXyeTsif18fKy6WOiWeyJxGyVFnOQAw9WWiLMd8WhklAcHfuX6jLIH3MwaxpUPV7qVdqOqRQ",
	"QuaSjk77ZE8/SCqkwDYK+1aiDIfB65jzd8xfvzztSeSXP/5hBsuXx/PcfHPu5y7AOsEXgISlKamSSYQi",
	"jDwLTEjGYqF8lI38b7PpXs/XZlmlgN5zlu18rv27s+xQxgvXy0kntyz1ENxssvCl/DWbbF5ruPpLa7Kn",
	"4dzbVLfKFpXuUtdsn11B9sN6pPnF+B9rk670fO+EHb8N98vng3rxY7cXf1hGvbdEG5Irr9j9/z1eb0cP",
	"b56q5NSidBSGqeYsTbTF0afhpyHdzDd/BwAA//+TDxdtvBMAAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
