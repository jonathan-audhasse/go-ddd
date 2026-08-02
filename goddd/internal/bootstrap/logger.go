package bootstrap

import (
	"goddd/internal/infrastructure/logger"

	"github.com/rs/zerolog"
)

func LoadLogger(level string) zerolog.Logger {
	return logger.New(level)
}
