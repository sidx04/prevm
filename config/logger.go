package config

import (
	"os"
	"time"

	"github.com/charmbracelet/log"
)

var Logger *log.Logger

func init() {
	Logger = log.NewWithOptions(os.Stderr, log.Options{
		ReportCaller:    false,
		ReportTimestamp: true,
		TimeFormat:      time.DateTime,
		Level:           log.DebugLevel,
		Prefix:          "EVM",
	})

	Logger.Info("Hello, Ethereum!")
}
