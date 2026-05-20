package main

import (
	"context"
	"log"

	httpapi "goddd/api/http"
	"goddd/internal/bootstrap"
)

func main() {

	app, err := bootstrap.New(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	api := httpapi.NewAPI(
		app.Config.ApiPort,
		app.Services,
		app.DB,
	)

	api.Serve()
}
