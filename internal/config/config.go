package config

import (
	"encoding/json"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName string `envconfig:"SERVICE_NAME" default:"shortener"`
	Server      struct {
		Listen      string        `envconfig:"SERVER_LISTEN_HTTP" default:":8080"`
		ReadTimeout time.Duration `envconfig:"SERVER_READ_TIMEOUT" default:"5s"`
		IdleTimeout time.Duration `envconfig:"SERVER_IDLE_TIMEOUT" default:"5s"`
	}
	Logger struct {
		Level  string `envconfig:"LOG_LEVEL" default:"info"`
		Output string `envconfig:"LOG_OUTPUT" default:"stdout"`
		Format string `envconfig:"LOG_FORMAT" default:"text"`
	}
}

// New parses environments and creates new instance of config
func New() (*Config, error) {
	cfg := new(Config)

	err := envconfig.Process("", cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) String() string {
	if out, err := json.MarshalIndent(&c, "", "  "); err == nil {
		return string(out)
	}
	return ""
}
