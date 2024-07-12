package memory

import (
	"goddd/src/domain/repository"
)

// Factory
func NewRepository() *repository.Repository {
	return &repository.Repository{CustomerRepo: NewCustomerRepo()}
}
