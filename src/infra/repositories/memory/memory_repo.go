package memory

// Memory repositories
type MemoryRepositories struct {
	CustRepo *CustomerMemoryRepository
}

// Factory
func NewMemoryRepositories() *MemoryRepositories {
	return &MemoryRepositories{CustRepo: NewCustomerMemoryRepository()}
}
