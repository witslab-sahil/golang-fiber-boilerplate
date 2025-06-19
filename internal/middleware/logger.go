package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/witslab-sahil/fiber-boilerplate/pkg/logger"
)

func Logger(log logger.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		path := c.Path()
		raw := string(c.Request().URI().QueryString())

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)
		
		// Get status code
		status := c.Response().StatusCode()

		// Get request ID
		requestID, _ := c.Locals("request_id").(string)

		// Log request details
		log.WithFields(map[string]interface{}{
			"request_id": requestID,
			"method":     c.Method(),
			"path":       path,
			"query":      raw,
			"ip":         c.IP(),
			"user_agent": c.Get("User-Agent"),
			"status":     status,
			"latency_ms": latency.Milliseconds(),
			"error":      err,
		}).Info("Request processed")

		return err
	}
}