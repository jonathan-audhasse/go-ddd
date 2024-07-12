package postgres

import (
	"goddd/src/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type healthRepo struct {
	db *sqlx.DB
}

func NewHealthRepo(db *sqlx.DB) *healthRepo {
	return &healthRepo{db}
}

func (r *healthRepo) HealthValue() (string, error) {
	var value string
	if err := r.db.Get(&value, "SELECT * FROM healthcheck"); err != nil {
		return "", errors.RepoItemNotFound.New("healthcheck value not found in repository")
	}
	return value, nil
}
