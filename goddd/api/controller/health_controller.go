package controller

import (
	"encoding/json"
	"net/http"

	"goddd/api/httperror"
	"goddd/internal/application/health"
)

type HealthController struct {
	service health.HealthService
}

func NewHealthController(service health.HealthService) *HealthController {
	return &HealthController{service: service}
}

func (c *HealthController) Health(w http.ResponseWriter, r *http.Request) {
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
