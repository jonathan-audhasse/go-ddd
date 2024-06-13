package services

import (
	"errors"
	"fmt"
	"goddd/src/domain/models"
	"goddd/src/services/dto"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CustomerRepository interface {
	Get(string) (models.Customer, error)
	Add(models.Customer) error
	List() ([]models.Customer, error)
}

type fakeCustomerRepo struct {
	mockGet  func(string) (models.Customer, error)
	mockAdd  func(models.Customer) error
	mockList func() ([]models.Customer, error)
}

func (fcr *fakeCustomerRepo) Get(id string) (models.Customer, error) {
	if fcr.mockGet != nil {
		return fcr.mockGet(id)
	}
	return models.Customer{}, nil
}

func (fcr *fakeCustomerRepo) Add(cust models.Customer) error {
	if fcr.mockAdd != nil {
		return fcr.mockAdd(cust)
	}
	return nil
}

func (fcr *fakeCustomerRepo) List() ([]models.Customer, error) {
	if fcr.mockList != nil {
		return fcr.mockList()
	}
	return []models.Customer{}, nil
}

// Test list customers
func TestCustomerServices(t *testing.T) {
	mokcErr := fmt.Errorf("mock error")
	repo := &fakeCustomerRepo{}
	// mokcErr := fmt.Errorf("mock error")
	services, _ := NewCustomerService(repo)
	// listing customers
	t.Run("list customers", func(t *testing.T) {
		t.Run("list customers should fail and wrap error", func(t *testing.T) {
			repo.mockList = func() ([]models.Customer, error) {
				return []models.Customer{}, mokcErr
			}
			_, err := services.ListCustomers()
			uerr := errors.Unwrap(err)
			assert.NotNil(t, uerr, "error should be wrapped")
			assert.ErrorIsf(t, mokcErr, uerr, "errors should be equals")
		})
		t.Run("list customers OK", func(t *testing.T) {
			repo.mockList = func() ([]models.Customer, error) {
				return []models.Customer{{}}, nil
			}
			res, err := services.ListCustomers()
			assert.Nil(t, err, "unexpected error")
			assert.Len(t, res, 1, "expected a list of customer")
		})
	})
	// getting customers
	t.Run("get a customer", func(t *testing.T) {
		t.Run("get a customer should fail and wrap error", func(t *testing.T) {
			repo.mockGet = func(id string) (models.Customer, error) {
				return models.Customer{}, mokcErr
			}
			_, err := services.GetCustomer("an_ID")
			uerr := errors.Unwrap(err)
			assert.NotNil(t, uerr, "error should be wrapped")
			assert.ErrorIsf(t, mokcErr, uerr, "errors should be equals")
		})
		t.Run("get customer OK", func(t *testing.T) {
			cust := models.Customer{ID: "an_ID", Name: "Johnny", Email: "a@b.c"}
			repo.mockGet = func(id string) (models.Customer, error) {
				return cust, nil
			}
			res, err := services.GetCustomer(cust.ID)
			assert.Nil(t, err, "unexpected error")
			assert.Equal(t, cust, res, "expected customer be the same")
		})
	})
	// adding customers
	t.Run("add a new customer", func(t *testing.T) {
		custToAdd := dto.CustomerDTO{Name: "Johnny", Email: "a@b.c"}
		t.Run("name should not be empty", func(t *testing.T) {
			_, err := services.AddCustomer(dto.CustomerDTO{})
			uerr := errors.Unwrap(err)
			assert.NotNil(t, uerr, "error should be wrapped")
			assert.ErrorIsf(t, models.ErrInvalidCustomer, uerr, "invalid customer should be return here")
		})
		t.Run("name should not be empty", func(t *testing.T) {
			repo.mockAdd = func(c models.Customer) error {
				return mokcErr
			}
			_, err := services.AddCustomer(custToAdd)
			uerr := errors.Unwrap(err)
			assert.NotNil(t, uerr, "error should be wrapped")
			assert.ErrorIsf(t, mokcErr, uerr, "errors should be equals")
		})
		t.Run("get customer OK", func(t *testing.T) {
			repo.mockAdd = func(c models.Customer) error {
				return nil
			}
			cust := models.Customer{Name: "Johnny", Email: "a@b.c"}
			res, err := services.AddCustomer(custToAdd)
			assert.Nil(t, err, "unexpected error")
			assert.Equal(t, cust.Name, res.Name, "expected customer name be the same")
			assert.Equal(t, cust.Email, res.Email, "expected customer email be the same")
		})
	})
}
