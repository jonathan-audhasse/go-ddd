package userservice_test

import (
	"context"
	mockRepo "goddd/tests/mocks/repository"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// FIXME use mock instead
func TestUserService_Create(t *testing.T) {
	assert.Fail(t, "to complete")
	ctx := context.Background()

	repo := mockRepo.NewMockUserRepository(t)
	repo.EXPECT().Create(ctx, mock.AnythingOfType("models.NewUser")).Return(gofakeit.Error())
}
