package api

import "os"

func getEnvDefault(key, def string) string {
	env := os.Getenv(key)
	if len(env) != 0 {
		return env
	}
	return def
}

var apiPort string = getEnvDefault("API_PORT", "8000")
