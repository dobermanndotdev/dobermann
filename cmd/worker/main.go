package main

import (
	"net/http"
	"strings"

	"github.com/kelseyhightower/envconfig"

	"github.com/flowck/dobermann/backend/internal/common/logs"
)

type Config struct {
	AmqpUrl     string `envconfig:"AMQP_URL"`
	Port        int    `envconfig:"HTTP_PORT"`
	DebugMode   string `envconfig:"DEBUG_MODE"`
	DatabaseURL string `envconfig:"DATABASE_URL"`
	Region      string `envconfig:"FLY_REGION" required:"true"`
}

func (c Config) IsDebugMode() bool {
	return strings.ToLower(c.DebugMode) == "enabled"
}

func main() {
	config := &Config{}
	err := envconfig.Process("", config)
	if err != nil {
		panic(err)
	}

	logger := logs.New(config.IsDebugMode())
	logger.Infof("Worker is running from region %s", config.Region)
	logger.Warn(http.ListenAndServe(":6060", nil))
}
