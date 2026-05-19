package controller

import (
	"encoding/json"
	"net/http"

	"goddd/api/httperror"
	healthservice "goddd/internal/application/health"
)

func Health(srv healthservice.HealthService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := srv.IsHealthy(r.Context()); err != nil {
			// not healthy
			httperror.Write(w, err)
			return
		}

		// healthy
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		_ = json.NewEncoder(w).Encode(map[string]string{})
	}
}
