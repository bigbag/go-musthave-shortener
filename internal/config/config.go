package config

import (
	"encoding/json"
	"flag"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	ServiceName      string `envconfig:"SERVICE_NAME" default:"shortener"`
	BaseURL          string `envconfig:"BASE_URL"`
	FileStoragePath  string `envconfig:"FILE_STORAGE_PATH"`
	UserCookieSecret string `envconfig:"USER_COOKIE_SECRET" default:"secret"`
	UserContextKey   string `envconfig:"USER_CONTEXT_KEY" default:"userid"`
	Server           struct {
		Listen      string        `envconfig:"SERVER_ADDRESS"  default:":8080"`
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

	flag.StringVar(&cfg.Server.Listen, "a", cfg.Server.Listen, "listen address. env: SERVER_ADDRESS")
	flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base url for short link. env: BASE_URL")
	flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path. env: FILE_STORAGE_PATH")
	flag.Parse()

	return cfg, nil
}

func (c *Config) String() string {
	if out, err := json.MarshalIndent(&c, "", "  "); err == nil {
		return string(out)
	}
	return ""
}
