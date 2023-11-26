package http

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/flowck/dobermann/backend/internal/app"
	"github.com/flowck/dobermann/backend/internal/common/logs"
)

type Port struct {
	server *http.Server
	config Config
}

type Config struct {
	Port        int
	Application *app.App
	Logger      *logs.Logger
	Ctx         context.Context
}

func NewPort(config Config) *Port {
	router := echo.New()
	router.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"message": "i am ok"})
	})

	return &Port{
		config: config,
		server: &http.Server{
			ReadTimeout:       time.Second * 30,
			ReadHeaderTimeout: time.Second * 30,
			WriteTimeout:      time.Second * 30,
			IdleTimeout:       time.Second * 30,
			Handler:           router,
			Addr:              fmt.Sprintf(":%d", config.Port),
			BaseContext: func(listener net.Listener) context.Context {
				return config.Ctx
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
