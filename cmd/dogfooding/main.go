package main

import (
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"

	"github.com/flowck/dobermann/backend/internal/common/logs"
	"github.com/flowck/dobermann/backend/tests/client"
)

const (
	maxMonitors      = 100
	delayPerPostInMs = 100
	columnWebsite    = 1
)

func main() {
	cli, err := client.NewClientWithResponses("https://api.dobermann.dev", client.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
		req.Header.Add(echo.HeaderAuthorization, fmt.Sprintf("Bearer %s", os.Getenv("DOGFOODING_JWT")))
		return nil
	}))
	if err != nil {
		logs.Fatal(err)
	}

	app := &App{cli: cli}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = loadDatasetInStreams(ctx, "cmd/dogfooding/Web_Scrapped_websites.csv", app.createMonitor)
	if err != nil {
		logs.Fatal(err)
	}

	logs.Info("End of execution")
}

type App struct {
	cli *client.ClientWithResponses
}

func (a *App) createMonitor(ctx context.Context, url string) error {
	logs.Infof("URL: %s", url)

	resp, err := a.cli.CreateMonitor(ctx, client.CreateMonitorRequest{
		CheckIntervalInSeconds: 60 * 3,
		EndpointUrl:            fmt.Sprintf("https://%s", strings.TrimPrefix(url, "https://")),
	})
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("failed to create monitor due to status code %d whereas %d is the expected status code", resp.StatusCode, http.StatusCreated)
	}

	return nil
}

func loadDatasetInStreams(ctx context.Context, fileName string, handler func(ctx context.Context, url string) error) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.Open(fmt.Sprintf("%s/%s", wd, fileName))
	if err != nil {
		return err
	}
	defer func() { _ = file.Close() }()

	counter := 0
	reader := csv.NewReader(file)
	for {
		line, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		if err != nil {
			return fmt.Errorf("an error occurred while reading the csv: %v", err)
		}

		if counter == 0 {
			counter++
			continue
		}

		err = handler(ctx, line[columnWebsite])
		if err != nil {
			return err
		}

		if counter == maxMonitors {
			break
		}

		counter++
		time.Sleep(time.Millisecond * time.Duration(delayPerPostInMs))
	}

	return nil
}
