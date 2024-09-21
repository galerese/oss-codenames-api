package configs

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

type Config struct {
	HttpPort     string `envconfig:"HTTP_PORT" required:"true"`
	LogLevel     string `envconfig:"LOG_LEVEL" required:"true" default:"info"`
	LogFormat    string `envconfig:"LOG_FORMAT" required:"true" default:"json"`
	DatabaseURL  string `envconfig:"DATABASE_URL" required:"true"`
	DatabaseName string `envconfig:"DATABASE_NAME" required:"true"`
	AppName      string `envconfig:"APP_NAME" required:"true"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := envconfig.Process("", &cfg)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load configuration from environment")
	}

	return &cfg, nil
}
