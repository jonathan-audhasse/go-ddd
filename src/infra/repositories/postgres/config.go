package postgres

import (
	"context"
	"fmt"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Host string `env:"DB_HOST, default=db"`
	Port int    `env:"DB_PORT, default=5432"`
	Name string `env:"DB_NAME, default=postgres"`
	User string `env:"DB_USER, default=admin"`
	Pwd  string `env:"DB_PWD, default=abc123"`
}

func NewConfig() Config {
	var cfg Config
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatal(err)
	}
	return cfg
}

func (cfg Config) Url() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Pwd, cfg.Name)
}
