package handler

import (
	"encoding/json"
	"net/http"

	"goddd/api/http/httperror"
	"goddd/internal/application/health"
)

// HealthHandler holds the health handler
type HealthHandler struct {
	service health.HealthService
}

// NewHealthHandler instantiates a new HealthHandler
func NewHealthHandler(service health.HealthService) *HealthHandler {
	return &HealthHandler{service: service}
}

// IsHealthy calls IsHealthy service
func (c *HealthHandler) IsHealthy(w http.ResponseWriter, r *http.Request) {
	if err := c.service.IsHealthy(r.Context()); err != nil {
		// not healthy
		httperror.Write(w, err)
		return
	}

	// healthy
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(map[string]string{})
}
