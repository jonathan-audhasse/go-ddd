package main

import (
	"context"
	"flag"

	"goddd/internal/bootstrap"
	"goddd/internal/infrastructure/persistence/migrations"
)

const schemaName = "goddd"

func main() {

	cmd := flag.String("cmd", "up", "up | down | steps | version")
	steps := flag.Int("steps", 1, "steps count")
	flag.Parse()

	app, err := bootstrap.New(context.Background())
	if err != nil {
		panic(err)
	}

	switch *cmd {

	case "up":
		err = migrations.Up(app.DB, schemaName)

	case "down":
		err = migrations.Down(app.DB, schemaName)

	case "steps":
		err = migrations.Steps(app.DB, schemaName, *steps)

	case "version":
		v, dirty, err := migrations.Version(app.DB, schemaName)

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
