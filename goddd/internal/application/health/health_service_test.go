package health_test

import (
	"context"
	"goddd/internal/application/health"
	mockRepo "goddd/tests/mocks/repository"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
)

func TestService_IsReady(t *testing.T) {
	ctx := context.Background()

	t.Run("not ready", func(t *testing.T) {
		// mock repo behavior
		repo := mockRepo.NewMockHealthRepository(t)
		repo.EXPECT().Ping(ctx).Return(gofakeit.Error())

		srv := health.NewService(repo)
		assert.False(t, srv.IsReady(ctx))
	})

	t.Run("ready", func(t *testing.T) {
		// mock repo behavior
		repo := mockRepo.NewMockHealthRepository(t)
		repo.EXPECT().Ping(ctx).Return(nil)

		srv := health.NewService(repo)
		assert.True(t, srv.IsReady(ctx))
	})
}
