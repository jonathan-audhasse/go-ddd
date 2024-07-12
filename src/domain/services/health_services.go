package services

import (
	"goddd/src/domain/repository"
	"goddd/src/pkg/errors"
	"log"
	"os"
)

// HealthService for asset's services implementation
type HealthService struct {
	repo repository.HealthRepository
}

// Create the HealthService
func NewHealthService(repo repository.HealthRepository) *HealthService {
	return &HealthService{repo}
}

// health check
func (hs *HealthService) IsHealthy() error {
	v, err := hs.repo.HealthValue()
	// unknow error
	if err != nil {
		log.Println(errors.FullError(err))
		return errors.Wrap(err, "healthcheck failure")
	}
	// check health
	if v != os.Getenv("HEALTH_VALUE") {
		return errors.InternalError.New("not healthy")
	}
	// healthy
	return nil
}
