package main

import (
	"flag"
	"goddd/db"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// --- Logging setup ---
	// Pretty console output in development; swap to zerolog.New(os.Stderr) for
	// JSON in production.
	log.Logger = log.Output(zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.RFC3339,
	})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// --- CLI flags (override with env vars in production) ---
	dsn := flag.String("dsn",
		getenv("DATABASE_URL", "postgres://admin:postgres@db:5432/myapp?sslmode=disable"),
		"Postgres DSN")
	migrationsDir := flag.String("migrations", "file://migrations", "Path to migration files")
	cmd := flag.String("cmd", "up", "Migration command: up | down | steps | version")
	steps := flag.Int("steps", 1, "Number of steps for the 'steps' command")
	flag.Parse()

	logger := log.With().Str("cmd", *cmd).Logger()

	// --- Open database ---
	sqlDB, err := db.Open(*dsn)
	if err != nil {
		logger.Fatal().Err(err).Msg("failed to open database")
	}
	defer sqlDB.Close()

	// --- Run requested migration command ---
	switch *cmd {
	case "up":
		if err := db.Up(sqlDB, *migrationsDir); err != nil {
			logger.Fatal().Err(err).Msg("migrate up failed")
		}
		logger.Info().Msg("all migrations applied")

	case "down":
		if err := db.Down(sqlDB, *migrationsDir); err != nil {
			logger.Fatal().Err(err).Msg("migrate down failed")
		}
		logger.Info().Msg("all migrations rolled back")

	case "steps":
		if err := db.Steps(sqlDB, *migrationsDir, *steps); err != nil {
			logger.Fatal().Err(err).Int("steps", *steps).Msg("migrate steps failed")
		}
		logger.Info().Int("steps", *steps).Msg("migration steps applied")

	case "version":
		v, dirty, err := db.Version(sqlDB, *migrationsDir)
		if err != nil {
			logger.Fatal().Err(err).Msg("could not read migration version")
		}
		logger.Info().Uint("version", v).Bool("dirty", dirty).Msg("current migration version")

	default:
		logger.Fatal().Str("got", *cmd).Msg("unknown command — use up, down, steps, or version")
	}
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}