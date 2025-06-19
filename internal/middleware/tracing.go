package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Tracing returns a Fiber middleware for OpenTelemetry tracing
func Tracing(serviceName string) fiber.Handler {
	tracer := otel.Tracer(serviceName)
	propagator := otel.GetTextMapPropagator()

	return func(c *fiber.Ctx) error {
		// Extract trace context from headers
		headers := make(map[string][]string)
		c.Request().Header.VisitAll(func(key, value []byte) {
			headers[string(key)] = []string{string(value)}
		})
		ctx := propagator.Extract(c.Context(), propagation.HeaderCarrier(headers))

		// Start span
		spanName := fmt.Sprintf("%s %s", c.Method(), c.Path())
		ctx, span := tracer.Start(ctx, spanName,
			trace.WithSpanKind(trace.SpanKindServer),
			trace.WithAttributes(
				attribute.String("http.method", c.Method()),
				attribute.String("http.target", c.Path()),
				attribute.String("http.url", c.OriginalURL()),
				attribute.String("http.user_agent", c.Get("User-Agent")),
				attribute.Int("http.request_content_length", c.Request().Header.ContentLength()),
				attribute.String("net.host.name", c.Hostname()),
			),
		)
		defer span.End()

		// Store context in Fiber locals
		c.Locals("otel-context", ctx)

		// Process request
		err := c.Next()

		// Set status code
		statusCode := c.Response().StatusCode()
		span.SetAttributes(attribute.Int("http.status_code", statusCode))

		// Record error if any
		if err != nil {
			span.RecordError(err)
			span.SetAttributes(attribute.Bool("error", true))
		}

		return err
	}
}