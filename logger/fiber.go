package logger

import (
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// FiberMiddleware creates a Fiber middleware using the provided logger
func FiberMiddleware(logger *zap.Logger) fiber.Handler {
	return fiberzap.New(fiberzap.Config{
		Logger: logger,
		Fields: []string{
			"latency", "status", "method", "url", "ip", "user_agent", "request_id",
			"bytes_sent", "bytes_received", "protocol", "host", "path", "query",
			"referer", "content_type", "content_length",
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
