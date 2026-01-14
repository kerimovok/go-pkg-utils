package logger

// Config holds logging configuration
type Config struct {
	Enabled    *bool  `yaml:"enabled"`
	FilePath   string `yaml:"file_path"`   // Path to log file (empty for stdout only)
	MaxSize    int64  `yaml:"max_size"`    // Max size in bytes before rotation
	MaxBackups int    `yaml:"max_backups"` // Max number of backup files to retain
	MaxAge     int    `yaml:"max_age"`     // Max age of backup files in days
	Level      string `yaml:"level"`       // Log level: debug, info, warn, error (default: info)
}

// IsEnabled returns true if logging is enabled
func (c *Config) IsEnabled() bool {
	if c == nil || c.Enabled == nil {
		return false
	}
	return *c.Enabled
}

// DefaultConfig returns a default logging configuration
func DefaultConfig() *Config {
	enabled := true
	return &Config{
		Enabled:    &enabled,
		FilePath:   "",
		MaxSize:    104857600, // 100MB
		MaxBackups: 3,
		MaxAge:     28,
		Level:      "info",
	}
}
