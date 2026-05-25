package main

import (
	"context"
	"flag"

	"goddd/internal/bootstrap"
	"goddd/internal/infrastructure/persistence/migrations"
)

func main() {

	cmd := flag.String("cmd", "up", "up | down | steps | version")
	steps := flag.Int("steps", 1, "steps count")
	flag.Parse()

	app, err := bootstrap.New(context.Background())
	if err != nil {
		panic(err)
	}

	dsn := app.Config.DatabaseURL()
	schemaName := app.SchemaName

	switch *cmd {

	case "up":
		err = migrations.Up(dsn, app.SchemaName)

	case "down":
		err = migrations.Down(dsn, schemaName)

	case "steps":
		err = migrations.Steps(dsn, schemaName, *steps)

	case "version":
		v, dirty, err := migrations.Version(dsn, schemaName)

		if err != nil {
			app.Logger.Fatal().Err(err).Msg("failed to get migration version")
		}

		app.Logger.Info().
			Uint("version", v).
			Bool("dirty", dirty).
			Msg("migration version")

		return
	}

	if err != nil {
		app.Logger.Fatal().
			Err(err).
			Str("cmd", *cmd).
			Msg("migration failed")
	}
}
