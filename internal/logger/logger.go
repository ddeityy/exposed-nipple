package logger

import (
	"os"

	"nipple/internal/config"

	"github.com/charmbracelet/log"
)

func New(cfg config.Logger) *log.Logger {
	var logLevel log.Level
	switch cfg.LogLevel {
	case "debug":
		logLevel = log.DebugLevel
	case "info":
		logLevel = log.InfoLevel
	case "warn":
		logLevel = log.WarnLevel
	case "error":
		logLevel = log.ErrorLevel
	default:
		logLevel = log.InfoLevel
	}

	return log.NewWithOptions(os.Stderr, log.Options{
		Level:           logLevel,
		ReportTimestamp: cfg.ReportTimestamp,
	})
}
