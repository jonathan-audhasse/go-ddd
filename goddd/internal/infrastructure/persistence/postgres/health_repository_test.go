package postgres_test

import (
	"context"
	"goddd/internal/infrastructure/persistence/postgres"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestHealthRepository_Ping(t *testing.T) {
	initTestLogger()

	ctx := log.Logger.WithContext(context.Background())

	// set repo
	repo := postgres.NewHealthRepository(testDB)
	// ping
	require.NoError(t, repo.Ping(ctx))
}
