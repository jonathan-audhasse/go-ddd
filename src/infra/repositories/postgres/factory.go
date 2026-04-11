package postgres

import (
	"goddd/domain/repository"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Factory
func NewRepository() (*repository.Repository, func(), error) {
	db, dbClose, err := Connect()
	if err != nil {
		return nil, nil, err
	}
	return &repository.Repository{
		UserRepo:     NewUserRepo(db),
		CustomerRepo: NewCustomerRepo(db),
	}, dbClose, nil
}

// connect to repository
func Connect() (*sqlx.DB, func(), error) {
	log.Println("Set up DB. Connecting...")
	cfg := NewConfig()
	db, err := sqlx.Connect("postgres", cfg.Url())
	if err != nil {
		log.Println("failed to connect to DB: ", err)
		return nil, nil, err
	}
	return db, func() {
		if err := db.Close(); err != nil {
			log.Panicln("failed to close DB session:", err)
		}
		log.Println("DB session closed")
	}, nil
}
