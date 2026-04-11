package postgres

import (
	"goddd/domain/models"
	"goddd/pkg/errors"
	"log"

	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db}
}

// Get a user by ID
func (r *userRepo) GetById(id string) (models.User, error) {
	if id == "" {
		return models.User{}, errors.InternalError.New("user id must not be empty")
	}
	var u models.User
	if err := r.db.Get(&u, "SELECT * FROM \"user\" WHERE id=$1", id); err != nil {
		log.Printf("failed to retrieve user by user id=%s\n", id)
		return models.User{}, errors.RepoItemNotFound.New("user not found in repository")
	}
	return u, nil
}

// Get a user by Username
func (r *userRepo) GetByUsername(username string) (models.User, error) {
	if username == "" {
		return models.User{}, errors.InternalError.New("user username must not be empty")
	}
	var u models.User
	if err := r.db.Get(&u, "SELECT * FROM \"user\" WHERE username=$1", username); err != nil {
		log.Printf("failed to retrieve user by username=%s: %s\n", username, err)
		return models.User{}, errors.RepoItemNotFound.New("user not found in repository")
	}
	return u, nil
}
