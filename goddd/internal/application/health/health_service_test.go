package health_test

import (
	"context"
	healthservice "goddd/internal/application/health"
	mockService "goddd/tests/mocks/service"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestHealthService_IsReady(t *testing.T) {
	ctx := context.Background()

	t.Run("not healthy", func(t *testing.T) {
		// mock repo behavior
		repo := mockService.NewMockPinger(t)
		repo.EXPECT().Ping(ctx).Return(gofakeit.Error())

		srv := healthservice.NewHealthService(repo)
		assert.Error(t, srv.IsHealthy(ctx))
	})

	t.Run("healthy", func(t *testing.T) {
		// mock repo behavior
		repo := mockService.NewMockPinger(t)
		repo.EXPECT().Ping(ctx).Return(nil)

		srv := healthservice.NewHealthService(repo)
		assert.NoError(t, srv.IsHealthy(ctx))
	})
}
