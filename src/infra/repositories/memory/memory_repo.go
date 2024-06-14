package memory

import (
	"goddd/src/domain/repository"
)

// Factory
func NewMemoryRepository() *repository.Repository {
	return &repository.Repository{CustomerRepo: NewCustomerMemoryRepository()}
}
