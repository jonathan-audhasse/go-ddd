package postgres_test

import (
	"context"
	"goddd/internal/infrastructure/persistence/postgres"
	"testing"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
)

func TestHealthPinger_Ping(t *testing.T) {
	ctx := log.Logger.WithContext(context.Background())

	// set repo
	repo := postgres.NewPinger(testDB)
	// ping
	require.NoError(t, repo.Ping(ctx))
}
