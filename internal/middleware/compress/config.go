package compress

import (
	"github.com/gofiber/fiber/v2"
)

type Level int

const (
	LevelNoCompression      Level = 0
	LevelBestSpeed          Level = 1
	LevelBestCompression    Level = 9
	LevelDefaultCompression Level = -1
)

type Config struct {
	Next       func(c *fiber.Ctx) bool
	Header     string
	BadRequest fiber.Handler
	Level      Level
}

// ConfigDefault is the default config
var ConfigDefault = Config{
	Next:       nil,
	Header:     fiber.HeaderAcceptEncoding,
	BadRequest: nil,
	Level:      LevelDefaultCompression,
}

// Helper function to set default values
func configDefault(config ...Config) Config {
	// Return default config if nothing provided
	if len(config) < 1 {
		return ConfigDefault
	}

	// Override default config
	cfg := config[0]

	// Set default values
	if cfg.Header == "" {
		cfg.Header = ConfigDefault.Header
	}

	if cfg.BadRequest == nil {
		cfg.BadRequest = func(c *fiber.Ctx) error {
			return c.SendStatus(fiber.StatusBadRequest)
		}
	}

	if cfg.Level < LevelNoCompression || cfg.Level > LevelBestCompression {
		cfg.Level = ConfigDefault.Level
	}

	return cfg
}
