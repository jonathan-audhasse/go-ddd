package bootstrap

import (
	"context"

	"goddd/internal/infrastructure/config"
	"goddd/internal/infrastructure/persistence/postgres"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewDatabase(ctx context.Context, cfg *config.Config, schemaName string) (*pgxpool.Pool, error) {
	return postgres.Open(ctx, cfg.DatabaseURL(), schemaName)
}
