package logger

import (
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// New configures the global zerolog logger and returns it.
// level should be one of: trace, debug, info, warn, error, fatal, panic.
func New(level string) zerolog.Logger {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)
 
	// Pretty console output — swap to zerolog.New(os.Stderr) for JSON in prod.
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	})
 
	return log.Logger
}
 