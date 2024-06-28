package api

import (
	"context"
	"fmt"
	"log"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	ApiPort int `env:"API_PORT, default=8000"`
	DB      *DbConfig
}

type DbConfig struct {
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

func (cfg Config) DbUrl() string {
	db := cfg.DB
	// return fmt.Sprintf("postgresql://%s:%s@%s:%d/%s", db.User, db.Pwd, db.Name, db.Port, db.Name)
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.Host, db.Port, db.User, db.Pwd, db.Name)
}
