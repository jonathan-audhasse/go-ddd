package api

import (
	"context"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	ApiPort int `env:"API_PORT, default=8000"`
}

func NewConfig() Config {
	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}
