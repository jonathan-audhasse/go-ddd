package config

import (
	"fmt"

	// File source driver — reads .sql files from disk
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/ilyakaznacheev/cleanenv"
)

// Config holds every environment variable the application needs.
// All values are injected by Docker Compose (or any other env provider).
type Config struct {
	// Database
	DBConfig

	// API HTTP server
	ApiPort int `env:"API_PORT" env-default:"8000"`

	// FIXME change to JWT
	// // Basic auth credentials (used by the basic_auth middleware)
	// BasicAuthUser     string `env:"BASIC_AUTH_USER"     env-required:"true"`
	// BasicAuthPassword string `env:"BASIC_AUTH_PASSWORD" env-required:"true"`

	// Logging
	LogLevel string `env:"LOG_LEVEL" env-default:"info"`
}

// DBConfig holds database configuration
type DBConfig struct {
	Host     string `env:"DB_HOST" env-default:"db"`
	Port     string `env:"DB_PORT" env-default:"5432"`
	Name     string `env:"DB_NAME" env-default:"postgres"`
	User     string `env:"DB_USER" env-default:"admin"`
	Password string `env:"DB_PWD" env-default:"abc123"`
	SslMode  string `env:"SSL_MODE" env-default:"prefer"`
}

// DatabaseURL return the database URL
// e.g. "postgres://user:password@localhost:5432/mydb?sslmode=prefer"
func (c DBConfig) DatabaseURL() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.Name, c.SslMode)
}

// Load reads environment variables into Config and returns an error
// if any required variable is missing.
func Load() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, fmt.Errorf("config.Load: %w", err)
	}
	return &cfg, nil
}
