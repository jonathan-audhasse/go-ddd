package services

import (
	"fmt"
	"goddd/domain/models"
	mockrepo "goddd/domain/repository/.mocks"
	"goddd/domain/services/dto"
	"goddd/pkg/errors"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// randomString generates a random string of length n
func randomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	k := len(charset)

	res := make([]byte, n)
	for i := range res {
		res[i] = charset[rand.Intn(k)]
	}
	return string(res)
}

// Test ListCustomers
func TestCustomerServices_List(t *testing.T) {
	anId, mokcErr := randomString(3), fmt.Errorf("mock error")
	testCases := []struct {
		name   string
		userId string
		expErr error
	}{
		{
			name:   "list customers: errors should be catched and wrapped",
			expErr: errors.Wrap(mokcErr, "fail to list customers"),
		},
		{
			name:   "list customers OK",
			userId: anId,
			expErr: nil,
		},
	}

	custumerRepo := mockrepo.NewMockCustomerRepository(t)
	// set the mock customer repo ListByUser behaviour
	custumerRepo.EXPECT().ListByUser(mock.AnythingOfType("string")).RunAndReturn(func(userid string) ([]models.Customer, error) {
		if userid != anId {
			// fail with an error
			return []models.Customer{}, mokcErr
		}
		// return singleton
		return []models.Customer{{}}, nil
	})

	// define the service
	services := NewCustomerService(custumerRepo)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// run service
			res, err := services.ListCustomers(tc.userId)
			if tc.expErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expErr, err)
				return
			}
			assert.NoError(t, err, "unexpected error")
			assert.Len(t, res, 1)
		})
	}
}

// Test GetCustomer
func TestCustomerServices_Get(t *testing.T) {
	cust := models.Customer{Id: randomString(3), Name: "Johnny", Email: "a@b.c"}
	mokcErr := fmt.Errorf("mock error")
	custumerRepo := mockrepo.NewMockCustomerRepository(t)
	// set the mock customer repo Get behaviour
	custumerRepo.EXPECT().Get(mock.AnythingOfType("string")).RunAndReturn(func(id string) (models.Customer, error) {
		if id != cust.Id {
			// fail with an error
			return models.Customer{}, mokcErr
		}
		// return singleton
		return cust, nil
	})

	testCases := []struct {
		name   string
		id     string
		expErr error
		expRes models.Customer
	}{
		{
			name:   "get a customer: errors should be catched and wrapped",
			id:     "unknowId",
			expErr: errors.Wrap(mokcErr, "fail to get customer id=unknowId"),
		},
		{
			name:   "get customer OK",
			id:     cust.Id,
			expErr: nil,
			expRes: cust,
		},
	}

	// define the service
	services := NewCustomerService(custumerRepo)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// run service
			res, err := services.GetCustomer(tc.id)
			if tc.expErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expErr, err)
				return
			}
			assert.NoError(t, err, "unexpected error")
			assert.Equal(t, tc.expRes, res, "expected customer be the same")
		})
	}
}

// Test DeleteCustomer
func TestCustomerServices_Delete(t *testing.T) {
	anId, mokcErr := randomString(3), fmt.Errorf("mock error")
	custumerRepo := mockrepo.NewMockCustomerRepository(t)
	// set the mock customer repo Delete behaviour
	custumerRepo.EXPECT().Delete(mock.AnythingOfType("string")).RunAndReturn(func(id string) error {
		if id != anId {
			// fail with an error
			return mokcErr
		}
		// no erro
		return nil
	})

	testCases := []struct {
		name   string
		id     string
		expErr error
	}{
		{
			name:   "delete a customer: errors should be catched and wrapped",
			id:     "unknowId",
			expErr: errors.Wrap(mokcErr, "fail to delete customer id=unknowId"),
		},
		{
			name:   "delete customer OK",
			id:     anId,
			expErr: nil,
		},
	}

	// define the service
	services := NewCustomerService(custumerRepo)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// run service
			err := services.DeleteCustomer(tc.id)
			if tc.expErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expErr, err)
				return
			}
			assert.NoError(t, err, "unexpected error")
		})
	}
}

// Test AddCustomer
func TestCustomerServices_Add(t *testing.T) {
	nameTriggerErro, mokcErr := "nameTriggerErro", fmt.Errorf("mock error")
	custumerRepo := mockrepo.NewMockCustomerRepository(t)
	// set the mock customer repo Get behaviour
	custumerRepo.EXPECT().Add(mock.AnythingOfType("models.Customer")).RunAndReturn(func(cust models.Customer) error {
		if cust.Name == "nameTriggerErro" {
			// fail with an error
			return mokcErr
		}
		// no error
		return nil
	}).Maybe()

	testCases := []struct {
		name   string
		cust   dto.CustomerDTO
		expErr error
	}{
		{
			name:   "add customer: name should not be empty",
			cust:   dto.CustomerDTO{},
			expErr: errors.Wrap(errors.InvalidFormat.New("customer name should not be empty"), "fail to add customer"),
		},
		{
			name:   "add customer: errors should be catched and wrapped",
			cust:   dto.CustomerDTO{Name: nameTriggerErro, Email: "a@b.c"},
			expErr: errors.Wrap(mokcErr, "fail to add customer ({nameTriggerErro a@b.c})"),
		},
		{
			name:   "add customer: errors should be catched and wrapped",
			cust:   dto.CustomerDTO{Name: "Johnny", Email: "a@b.c"},
			expErr: nil,
		},
	}

	// define the service
	services := NewCustomerService(custumerRepo)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// run service
			res, err := services.AddCustomer("userid", tc.cust)
			if tc.expErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tc.expErr, err)
				return
			}
			assert.NoError(t, err, "unexpected error")
			// assert customer id not null
			assert.Contains(t, res.Id, "cust_")
			assert.Equal(t, tc.cust.Name, res.Name, "expected customer name be the same")
			assert.Equal(t, tc.cust.Email, res.Email, "expected customer email be the same")
		})
	}
}
