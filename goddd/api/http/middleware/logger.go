package middleware

import (
	"net/http"

	"github.com/rs/zerolog/log"
)

// LoggerHandler handles all context logs
func LoggerHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := log.Logger.WithContext(r.Context())
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
