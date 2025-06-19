package middleware

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

// ErrorHandler is the custom error handler for Fiber
func ErrorHandler(c *fiber.Ctx, err error) error {
	// Default to 500 server error
	code := fiber.StatusInternalServerError
	message := "Internal Server Error"

	// Check if it's a Fiber error
	var fiberErr *fiber.Error
	if errors.As(err, &fiberErr) {
		code = fiberErr.Code
		message = fiberErr.Message
	}

	// Custom error types can be handled here
	switch {
	case errors.Is(err, fiber.ErrUnauthorized):
		code = fiber.StatusUnauthorized
		message = "Unauthorized"
	case errors.Is(err, fiber.ErrForbidden):
		code = fiber.StatusForbidden
		message = "Forbidden"
	case errors.Is(err, fiber.ErrNotFound):
		code = fiber.StatusNotFound
		message = "Not Found"
	case errors.Is(err, fiber.ErrBadRequest):
		code = fiber.StatusBadRequest
		message = "Bad Request"
	}

	// Return JSON error response
	return c.Status(code).JSON(fiber.Map{
		"error": message,
	})
}