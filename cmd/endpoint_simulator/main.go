package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/flowck/dobermann/backend/internal/common/logs"
)

var failureStatusCode = []int{
	http.StatusOK,
	http.StatusBadRequest,
	http.StatusForbidden,
	http.StatusTemporaryRedirect,
	http.StatusInternalServerError,
}

func getStatusCode() int {
	return failureStatusCode[gofakeit.Number(0, len(failureStatusCode)-1)]
}

func main() {
	gofakeit.Seed(0)
	logger := logs.New(true)

	router := echo.New()
	router.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogMethod: true,
		LogError:  false,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			fields := logs.Fields{
				"method": v.Method,
				"uri":    v.URI,

				"status": v.Status,
			}

			if v.Error != nil || v.Status > http.StatusBadRequest {
				logger.WithFields(fields).Error("request handled with an error")
			} else {
				logger.WithFields(fields).Info("request handled successfully")
			}

			return nil
		},
	}))
	router.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Add("Server", "EndpointSimulator")
			return next(c)
		}
	})
	router.Use(middleware.Recover())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	router.GET("/", mainHandler)
	router.GET("/random", randomFailureHandler)

	port := "8090"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	server := http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadTimeout:       time.Minute * 2,
		ReadHeaderTimeout: time.Minute * 2,
		WriteTimeout:      time.Minute * 2,
		IdleTimeout:       time.Minute * 2,
		BaseContext: func(listener net.Listener) context.Context {
			return ctx
		},
	}

	go func() {
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Println(err)
			return
		}
	}()

	<-done
	tctx, tcancel := context.WithTimeout(context.Background(), time.Second*2)
	defer tcancel()
	err := server.Shutdown(tctx)
	if err != nil {
		os.Exit(1)
	}
}
