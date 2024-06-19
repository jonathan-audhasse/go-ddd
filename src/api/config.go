package api

import (
	"os"
)

type Config struct {
	Port string `envconfig:"PORT" default:"8080"`
}

func getEnvDefault(key, def string) string {
	env := os.Getenv(key)
	if len(env) != 0 {
		return env
	}
	return def
}

func NewConfig() Config {
	cfg := Config{
		Port: getEnvDefault("API_PORT", "8000"),
	}

	// if err := envconfig.Process("", &cfg); err != nil {
	// 	log.Fatalf("parseConfigError: %v", err)
	// }
	return cfg
}
