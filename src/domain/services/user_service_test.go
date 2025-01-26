package services

import (
	"fmt"
	"goddd/src/domain/models"
	"goddd/src/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type UserRepository interface {
	GetById(string) (models.User, error)
	GetByUsername(string) (models.User, error)
}

type fakeUserRepo struct {
	mockGetById       func(string) (models.User, error)
	mockGetByUsername func(string) (models.User, error)
}

func (fr *fakeUserRepo) GetById(id string) (models.User, error) {
	if fr.mockGetById != nil {
		return fr.mockGetById(id)
	}
	return models.User{}, nil
}

func (fr *fakeUserRepo) GetByUsername(uname string) (models.User, error) {
	if fr.mockGetByUsername != nil {
		return fr.mockGetByUsername(uname)
	}
	return models.User{}, nil
}

// FIXME use mock instead
// Test GetUserByUsername
func TestUserServices_GetByUsername(t *testing.T) {
	mokcErr := fmt.Errorf("mock error")
	repo := &fakeUserRepo{}
	services := NewUserService(repo)
	t.Run("get a user: errors should be catched and wrapped", func(t *testing.T) {
		repo.mockGetByUsername = func(id string) (models.User, error) {
			return models.User{}, mokcErr
		}
		_, err := services.GetUserByUsername("a_username")
		assert.Equal(t, errors.Wrap(mokcErr, "failed to retrieve user"), err)
	})
	t.Run("get user OK", func(t *testing.T) {
		u := models.User{Id: "an_ID", Username: "Jonathan", Password: "jonathan_password"}
		repo.mockGetByUsername = func(id string) (models.User, error) {
			return u, nil
		}
		res, err := services.GetUserByUsername(u.Id)
		assert.Nil(t, err, "unexpected error")
		assert.Equal(t, u, res, "expected user be the same")
	})
}
