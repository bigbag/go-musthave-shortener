package compress

import (
	"bytes"
	"compress/gzip"
	"github.com/gofiber/fiber/v2"
	"strings"
)

func copyBytes(b []byte) []byte {
	tmp := make([]byte, len(b))
	copy(tmp, b)
	return tmp
}

func compress(body []byte) ([]byte, error) {
	var buf bytes.Buffer
	gz, err := gzip.NewWriterLevel(&buf, gzip.BestCompression)
	if err != nil {
		return nil, err
	}

	gz.Write(body)
	defer gz.Close()
	return buf.Bytes(), nil
}

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

		if !strings.Contains(c.Get(cfg.Header, ""), "gzip") {
			return c.Next()
		}

		// Continue stack
		if err := c.Next(); err != nil {
			return err
		}

		compressBody, err := compress(copyBytes(c.Response().Body()))
		if err != nil {
			return cfg.BadRequest(c)
		}
		c.Response().Header.SetContentLength(-1)
		c.Response().SetBodyRaw(compressBody)
		c.Set(fiber.HeaderContentEncoding, "gzip")

		return nil
	}
}
