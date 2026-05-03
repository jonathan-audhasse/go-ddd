package main

import (
	"context"
	"flag"

	"goddd/internal/infrastructure/config"
	"goddd/internal/infrastructure/logger"
	"goddd/internal/infrastructure/persistence/migrations"
	"goddd/internal/infrastructure/persistence/postgres"
)

func main() {
	cmd := flag.String("cmd", "up", "Migration command: up | down | steps | version")
	steps := flag.Int("steps", 1, "Number of steps (used with -cmd=steps)")
	flag.Parse()

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	log := logger.New(cfg.LogLevel)

	pool, err := postgres.Open(ctx, cfg.DatabaseURL())
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open database")
	}
	defer pool.Close()

	switch *cmd {
	case "up":
		err = migrations.Up(pool, cfg.MigrationsPath)
	case "down":
		err = migrations.Down(pool, cfg.MigrationsPath)
	case "steps":
		err = migrations.Steps(pool, cfg.MigrationsPath, *steps)
	case "version":
		v, dirty, verr := migrations.Version(pool, cfg.MigrationsPath)
		if verr != nil {
			log.Fatal().Err(verr).Msg("could not read version")
		}
		log.Info().Uint("version", v).Bool("dirty", dirty).Msg("current migration version")
		return
	default:
		log.Fatal().Str("cmd", *cmd).Msg("unknown command — use up, down, steps, or version")
	}

	if err != nil {
		log.Fatal().Err(err).Str("cmd", *cmd).Msg("migration failed")
	}
}
