package repository

type Repository struct {
	HealthRepo HealthRepository
	UserRepo   UserRepository
}
