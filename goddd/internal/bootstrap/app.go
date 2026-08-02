package bootstrap

import (
	"context"

	"goddd/internal/application"
	"goddd/internal/infrastructure/config"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

// App holds the App bootstrap object
type App struct {
	Config     *config.Config
	Logger     zerolog.Logger
	DB         *pgxpool.Pool
	Services   *application.Services
	SchemaName string
}

// New instantiates a new app
func New(ctx context.Context) (*App, error) {
	// load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// load logger
	log := LoadLogger(cfg.LogLevel)

	// load database
	db, err := NewDatabase(ctx, cfg, cfg.SchemaName)
	if err != nil {
		return nil, err
	}

	// load repositories
	repos := LoadRepositories(db)

	// set services
	services := LoadServices(db, repos)

	return &App{
		Config:     cfg,
		Logger:     log,
		DB:         db,
		Services:   services,
		SchemaName: cfg.SchemaName,
	}, nil
}
