package services

import (
	"fmt"
	"goddd/src/domain/models"
	"goddd/src/domain/services/dto"
	"goddd/src/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CustomerRepository interface {
	Get(string) (models.Customer, error)
	Add(models.Customer) error
	ListByUser(userid string) ([]models.Customer, error)
	Delete(string) error
}

type fakeCustomerRepo struct {
	mockGet        func(string) (models.Customer, error)
	mockAdd        func(models.Customer) error
	mockListByUser func(userid string) ([]models.Customer, error)
	mockDelete     func(string) error
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

func (fcr *fakeCustomerRepo) ListByUser(userid string) ([]models.Customer, error) {
	if fcr.mockListByUser != nil {
		return fcr.mockListByUser(userid)
	}
	return []models.Customer{}, nil
}

func (fcr *fakeCustomerRepo) Delete(id string) error {
	if fcr.mockDelete != nil {
		return fcr.mockDelete(id)
	}
	return nil
}

// Test list customers
func TestCustomerServices_List(t *testing.T) {
	mokcErr := fmt.Errorf("mock error")
	repo := &fakeCustomerRepo{}
	services := NewCustomerService(repo)
	t.Run("list customers: errors should be catched and wrapped", func(t *testing.T) {
		repo.mockListByUser = func(string) ([]models.Customer, error) {
			return []models.Customer{}, mokcErr
		}
		_, err := services.ListCustomers("")
		assert.Equal(t, errors.Wrap(mokcErr, "fail to list customers"), err)
	})
	t.Run("list customers OK", func(t *testing.T) {
		repo.mockListByUser = func(string) ([]models.Customer, error) {
			return []models.Customer{{}}, nil
		}
		res, err := services.ListCustomers("userid")
		assert.Nil(t, err, "unexpected error")
		assert.Len(t, res, 1, "expected a list of customer")
	})
}

// Test get customer
func TestCustomerServices_Get(t *testing.T) {
	mokcErr := fmt.Errorf("mock error")
	repo := &fakeCustomerRepo{}
	services := NewCustomerService(repo)
	t.Run("get a customer: errors should be catched and wrapped", func(t *testing.T) {
		repo.mockGet = func(id string) (models.Customer, error) {
			return models.Customer{}, mokcErr
		}
		_, err := services.GetCustomer("an_ID")
		assert.Equal(t, errors.Wrap(mokcErr, "fail to get customer id=an_ID"), err)
	})
	t.Run("get customer OK", func(t *testing.T) {
		cust := models.Customer{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}
		repo.mockGet = func(id string) (models.Customer, error) {
			return cust, nil
		}
		res, err := services.GetCustomer(cust.Id)
		assert.Nil(t, err, "unexpected error")
		assert.Equal(t, cust, res, "expected customer be the same")
	})
}

// Test add customer
func TestCustomerServices_Add(t *testing.T) {
	mokcErr := fmt.Errorf("mock error")
	repo := &fakeCustomerRepo{}
	services := NewCustomerService(repo)
	// adding customers
	custToAdd := dto.CustomerDTO{Name: "Johnny", Email: "a@b.c"}
	t.Run("add customer: name should not be empty", func(t *testing.T) {
		_, err := services.AddCustomer("userid", dto.CustomerDTO{})
		assert.Equal(t, errors.GetID(err), errors.InvalidFormat)
		assert.EqualError(t, err, "fail to add customer")
		uerr := errors.Unwrap(err)
		assert.Equal(t, errors.InvalidFormat.New("customer name should not be empty"), uerr, "invalid customer should be return here")
	})
	t.Run("add customer: errors should be catched and wrapped", func(t *testing.T) {
		repo.mockAdd = func(c models.Customer) error {
			return mokcErr
		}
		_, err := services.AddCustomer("", custToAdd)
		assert.Equal(t, errors.Wrap(mokcErr, "fail to add customer ({Johnny a@b.c})"), err)
	})
	t.Run("add customer OK", func(t *testing.T) {
		repo.mockAdd = func(c models.Customer) error {
			return nil
		}
		cust := models.Customer{Name: "Johnny", Email: "a@b.c"}
		res, err := services.AddCustomer("userid", custToAdd)
		assert.Nil(t, err, "unexpected error")
		assert.Equal(t, cust.Name, res.Name, "expected customer name be the same")
		assert.Equal(t, cust.Email, res.Email, "expected customer email be the same")
	})
}

// Test delete customer
func TestCustomerServices_Delete(t *testing.T) {
	mokcErr := fmt.Errorf("mock error")
	repo := &fakeCustomerRepo{}
	services := NewCustomerService(repo)
	t.Run("delete a customer: errors should be catched and wrapped", func(t *testing.T) {
		repo.mockDelete = func(string) error {
			return mokcErr
		}
		assert.Equal(t, errors.Wrap(mokcErr, "fail to delete customer id=an_ID"), services.DeleteCustomer("an_ID"))
	})
	t.Run("delete customer OK", func(t *testing.T) {
		cust := models.Customer{Id: "an_ID"}
		repo.mockDelete = func(string) error {
			return nil
		}
		assert.NoError(t, services.DeleteCustomer(cust.Id), "unexpected error")
	})
}
