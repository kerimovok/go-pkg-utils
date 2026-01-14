package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// NewLogger creates a new Zap logger based on configuration
func NewLogger(config *Config) (*zap.Logger, error) {
	if config == nil || !config.IsEnabled() {
		// Return a no-op logger if logging is disabled
		return zap.NewNop(), nil
	}

	var logger *zap.Logger
	var err error

	// Configure Zap logger with Lumberjack for file rotation
	if config.FilePath != "" {
		// Production logger with file output
		writeSyncer := zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    int(config.MaxSize / (1024 * 1024)), // Convert bytes to MB
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   true,
		})

		// Also write to stdout in addition to file
		multiWriteSyncer := zapcore.NewMultiWriteSyncer(
			writeSyncer,
			zapcore.AddSync(os.Stdout),
		)

		core := zapcore.NewCore(
			zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
			multiWriteSyncer,
			parseLogLevel(config.Level),
		)
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	} else {
		// Development logger with console output
		devConfig := zap.NewDevelopmentConfig()
		devConfig.Level = zap.NewAtomicLevelAt(parseLogLevel(config.Level))
		logger, err = devConfig.Build()
		if err != nil {
			return nil, err
		}
	}

	return logger, nil
}

// parseLogLevel parses log level string to zapcore.Level
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// NewDevelopmentLogger creates a development logger (console output, colored)
func NewDevelopmentLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

// NewProductionLogger creates a production logger (JSON output, file rotation)
func NewProductionLogger(filePath string, maxSizeMB, maxBackups, maxAge int) (*zap.Logger, error) {
	config := &Config{
		Enabled:    func() *bool { b := true; return &b }(),
		FilePath:   filePath,
		MaxSize:    int64(maxSizeMB) * 1024 * 1024,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		Level:      "info",
	}
	return NewLogger(config)
}
