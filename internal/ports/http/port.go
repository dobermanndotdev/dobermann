package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimiddleware "github.com/oapi-codegen/echo-middleware"

	"github.com/flowck/dobermann/backend/internal/app"
	"github.com/flowck/dobermann/backend/internal/common/logs"
)

type Port struct {
	server *http.Server
	config Config
}

type Config struct {
	Port              int
	Application       *app.App
	AllowedCorsOrigin []string
	Logger            *logs.Logger
	JwtVerifier       jwtVerifier
	Ctx               context.Context
}

type handlers struct {
	application *app.App
}

func NewPort(config Config) (*Port, error) {
	router := echo.New()
	portHandlers := handlers{
		application: config.Application,
	}

	spec, err := GetSwagger()
	if err != nil {
		return nil, err
	}

	registerMiddlewares(router, spec, config)
	RegisterHandlers(router, portHandlers)

	router.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "i am ok"})
	})

	return &Port{
		config: config,
		server: &http.Server{
			Handler:           router,
			ReadTimeout:       time.Second * 30,
			ReadHeaderTimeout: time.Second * 30,
			WriteTimeout:      time.Second * 30,
			IdleTimeout:       time.Second * 30,
			Addr:              fmt.Sprintf(":%d", config.Port),
			BaseContext: func(listener net.Listener) context.Context {
				return config.Ctx
			},
		},
	}, nil
}

func (p *Port) Start() error {
	return p.server.ListenAndServe()
}

func (p *Port) Stop(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

func registerMiddlewares(router *echo.Echo, spec *openapi3.T, config Config) {
	router.HTTPErrorHandler = errorHandler(config.Logger)
	router.Use(middleware.RequestID())
	router.Use(correlationIdMiddleware)
	router.Use(middleware.Recover())
	router.Use(loggerMiddleware(config.Logger))
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: config.AllowedCorsOrigin,
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))

	spec.Servers = nil
	router.Use(oapimiddleware.OapiRequestValidatorWithOptions(spec, &oapimiddleware.Options{
		ErrorHandler: nil,
		Options: openapi3filter.Options{
			AuthenticationFunc: NewAuthenticator(config.JwtVerifier),
		},
	}))
}
