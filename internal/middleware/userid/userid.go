package userid

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// New creates a new middleware handler
func New(config ...Config) fiber.Handler {
	// Set default config
	cfg := configDefault(config...)

	// Return new handler
	return func(c *fiber.Ctx) error {
		// Don't execute middleware if Next returns true
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		storage := NewStorage(c, cfg.CookieName, cfg.Secret)

		userID, err := storage.Get()
		if err != nil {
			userID = uuid.New().String()
			storage.Set(userID)
		}

		// Add the user ID to locals
		c.Locals(cfg.ContextKey, userID)

		// // Continue stack
		return c.Next()
	}
}
