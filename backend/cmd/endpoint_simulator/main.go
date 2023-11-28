package main

import (
	"context"
	"errors"
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

func getRandomLatency() int {
	return gofakeit.Number(50, 250)
}

func main() {
	gofakeit.Seed(0)

	router := echo.New()
	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	router.GET("/", func(c echo.Context) error {
		if c.QueryParam("is_up") == "true" {
			return c.NoContent(http.StatusOK)
		}

		if c.QueryParam("is_up") == "false" {
			return c.NoContent(http.StatusInternalServerError)
		}

		// Quasi-random path
		s := getStatusCode()
		log.Println("Status code", s)
		time.Sleep(time.Millisecond * time.Duration(getRandomLatency()))
		return c.NoContent(s)
	})

	server := http.Server{
		Addr:              ":8090",
		Handler:           router,
		ReadTimeout:       time.Second * 30,
		ReadHeaderTimeout: time.Second * 30,
		WriteTimeout:      time.Second * 30,
		IdleTimeout:       time.Second * 30,
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
