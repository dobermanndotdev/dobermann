package main

import (
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/labstack/echo/v4"
)

const (
	minLatencyMs = 150
	maxLatencyMs = 300
)

func randomFailureHandler(c echo.Context) error {
	time.Sleep(time.Millisecond * time.Duration(gofakeit.Number(minLatencyMs, maxLatencyMs)))

	status := http.StatusOK

	if !gofakeit.Bool() {
		status = http.StatusInternalServerError
	}

	return c.JSON(status, map[string]string{
		"message": "ok",
	})
}

func mainHandler(c echo.Context) error {
	if c.QueryParam("timeout") == "true" {
		time.Sleep(time.Second * 10)
	}

	if c.QueryParam("is_up") == "true" {
		return c.NoContent(http.StatusOK)
	}

	if c.QueryParam("is_up") == "false" {
		return c.NoContent(http.StatusInternalServerError)
	}

	// Quasi-random path
	s := getStatusCode()
	return c.NoContent(s)
}
