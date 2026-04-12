package migrations

import (
	"errors"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/rs/zerolog/log"
)

// newMigrate bridges pgxpool → database/sql so go-migrate can use it.
// pgx/v5/stdlib.OpenDBFromPool wraps the pool without opening a new connection.
func newMigrate(pool *pgxpool.Pool, migrationsPath string) (*migrate.Migrate, error) {
	db := stdlib.OpenDBFromPool(pool)

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Err(err).Msg("failed to have a database driver")
		return nil, errors.New("failed to have a database driver")
	}
	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		log.Err(err).Msg("failed to have a new database driver instance")
		return nil, errors.New("failed to have a new database driver instance")
	}
	return m, nil
}

// Up applies every pending migration in ascending version order.
// Returns nil if already up to date.
func Up(pool *pgxpool.Pool, migrationsPath string) error {
	logger := log.With().Str("migrationsPath", migrationsPath).Logger()
	logger.Info().Msg("migrate database up...")
	m, err := newMigrate(pool, migrationsPath)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Err(err).Msg("failed to migrate DB schema up")
		return errors.New("failed to migrate DB schema up")
	}
	log.Info().Msg("migrates: up to date")
	return nil
}

// Down rolls back every applied migration, reverting the database to a clean state.
// Suitable for tests; avoid in production unless you know what you're doing.
func Down(pool *pgxpool.Pool, migrationsPath string) error {
	logger := log.With().Str("migrationsPath", migrationsPath).Logger()
	logger.Info().Msg("migrate database down...")
	m, err := newMigrate(pool, migrationsPath)
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Err(err).Msg("failed to migrate DB schema down")
		return errors.New("failed to migrate DB schema down")
	}
	log.Info().Msg("migrates: rolled back to zero")
	return nil
}

// Steps applies n migrations (positive = forward, negative = backward).
func Steps(pool *pgxpool.Pool, migrationsPath string, n int) error {
	logger := log.With().Str("migrationsPath", migrationsPath).Int("n", n).Logger()
	logger.Info().Msg("migrate database steps up ...")
	m, err := newMigrate(pool, migrationsPath)
	if err != nil {
		return err
	}
	if err := m.Steps(n); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Err(err).Msg("failed to migrate DB schema steps up")
		return errors.New("failed to migrate DB schema steps up")
	}
	return nil
}

// Version returns the currently applied migration version and whether the
// database is in a dirty state (a previous migration failed mid-way).
func Version(pool *pgxpool.Pool, migrationsPath string) (version uint, dirty bool, err error) {
	m, err := newMigrate(pool, migrationsPath)
	if err != nil {
		return 0, false, err
	}
	v, d, err := m.Version()
	if errors.Is(err, migrate.ErrNilVersion) {
		return 0, false, nil // no migration applied yet
	}
	return v, d, err
}