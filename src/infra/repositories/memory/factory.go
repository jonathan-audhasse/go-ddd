package memory

import (
	"goddd/domain/repository"
)

// Factory
func NewRepository() *repository.Repository {
	return &repository.Repository{CustomerRepo: NewCustomerRepo()}
}
