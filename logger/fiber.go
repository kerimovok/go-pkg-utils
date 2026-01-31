package logger

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// FiberMiddleware creates a Fiber middleware using the provided logger
// Stack traces are only included for 500-level (server) errors, not for client errors (4xx)
func FiberMiddleware(logger *zap.Logger) fiber.Handler {
	// Disable automatic stack traces - we'll add them conditionally for server errors only
	loggerWithoutStack := logger.WithOptions(zap.AddStacktrace(zapcore.FatalLevel))

	return fiberzap.New(fiberzap.Config{
		Logger: loggerWithoutStack,
		Fields: []string{
			"latency", "status", "method", "url", "ip", "user_agent", "request_id",
			"bytes_sent", "bytes_received", "protocol", "host", "path", "query",
			"referer", "content_type", "content_length",
		},
		FieldsFunc: func(c *fiber.Ctx) []zapcore.Field {
			statusCode := c.Response().StatusCode()

			// Only add stack trace for 500-level server errors
			// Client errors (4xx) don't need stack traces as they're expected business logic responses
			if statusCode >= 500 {
				return []zapcore.Field{zap.Stack("stacktrace")}
			}

			// No additional fields for client errors
			return []zapcore.Field{}
		},
	})
}

// SetupFiberLogger sets up structured logging for a Fiber app
// Returns the logger instance for cleanup (call logger.Sync() on shutdown)
func SetupFiberLogger(app *fiber.App, config *Config) (*zap.Logger, error) {
	if config == nil || !config.IsEnabled() {
		return zap.NewNop(), nil
	}

	logger, err := NewLogger(config)
	if err != nil {
		return nil, err
	}

	app.Use(FiberMiddleware(logger))
	return logger, nil
}
