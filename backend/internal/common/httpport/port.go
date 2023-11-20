package httpport

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/flowck/doberman/internal/common/logs"
)

type Port struct {
	server *http.Server
}

type PortConfig struct {
	Port              int
	AllowedCorsOrigin []string
	Router            *echo.Echo
	Logger            *logs.Logger
	Ctx               context.Context
}

func NewPort(cfg PortConfig) *Port {
	registerMiddlewares(cfg.Router, cfg)

	return &Port{
		server: &http.Server{
			Addr:              fmt.Sprintf(":%d", cfg.Port),
			Handler:           cfg.Router,
			ReadTimeout:       time.Second * 30,
			ReadHeaderTimeout: time.Second * 30,
			WriteTimeout:      time.Second * 30,
			IdleTimeout:       time.Second * 30,
			BaseContext: func(listener net.Listener) context.Context {
				return cfg.Ctx
			},
		},
	}
}

func (p *Port) Start() error {
	return p.server.ListenAndServe()
}

func (p *Port) Stop(ctx context.Context) error {
	return p.server.Shutdown(ctx)
}

func registerMiddlewares(router *echo.Echo, cfg PortConfig) {
	router.HTTPErrorHandler = errorHandler(cfg.Logger)
	router.Use(middleware.Recover())
	router.Use(middleware.RequestID())
	router.Use(loggerMiddleware(cfg.Logger))
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.AllowedCorsOrigin,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))
}
