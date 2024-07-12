package memory

import (
	"goddd/src/domain/models"
	"goddd/src/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerRepo_GetCustomer(t *testing.T) {
	// Create a fake customer to add to repository
	cust := models.Customer{Id: "an_ID", UserId: "", Name: "Johnny", Email: "a@b.c"}
	// Create the repo to use, and add some test Data to it for testing
	mr := customerRepo{
		customers: map[string]models.Customer{cust.Id: cust},
	}

	testCases := []struct {
		name   string
		id     string
		expErr error
	}{
		{
			name:   "Customer not found",
			id:     "another_ID",
			expErr: errors.RepoItemNotFound.New("customer id=another_ID not found in repository"),
		},
		{
			name:   "Customer Ok",
			id:     cust.Id,
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			_, err := mr.Get(tc.id)
			assert.Equal(t, err, tc.expErr)
		})
	}
}

func TestCustomerRepo_AddCustomer(t *testing.T) {
	mr := customerRepo{customers: map[string]models.Customer{}}

	cust := models.Customer{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}

	err := mr.Add(cust)
	if err != nil {
		t.Errorf("error not expected, got %v", err)
	}

	found, err := mr.Get(cust.Id)
	if err != nil {
		t.Fatal(err)
	}
	if found.Id != cust.Id {
		t.Errorf("expected %v, got %v", cust.Id, found.Id)
	}
}

func TestCustomerRepo_ListCustomers(t *testing.T) {
	// fake user to add to repository
	u := models.User{Id: "u_ID"}
	mr := customerRepo{customers: map[string]models.Customer{}}
	cust := models.Customer{Id: "an_ID", UserId: u.Id, Name: "Johnny", Email: "a@b.c"}

	// Test memory is empty
	ll, err := mr.ListByUser(cust.UserId)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, ll, 0, "repository should be empty")

	err = mr.Add(cust)
	if err != nil {
		t.Errorf("error not expected, got %v", err)
	}

	ll, err = mr.ListByUser(cust.UserId)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, ll, 1, "repository should have a single element")
}

func TestCustomerRepo_DeleteCustomer(t *testing.T) {
	// Create a fake customer to add to repository
	cust := models.Customer{Id: "an_ID", UserId: "userid", Name: "Johnny", Email: "a@b.c"}
	// Create the repo to use, and add some test Data to it for testing
	mr := customerRepo{
		customers: map[string]models.Customer{cust.Id: cust},
	}

	testCases := []struct {
		name   string
		id     string
		expErr error
		expLen int
	}{
		{
			name:   "Customer not found",
			id:     "another_ID",
			expErr: nil,
			expLen: 1,
		},
		{
			name:   "Customer Ok",
			id:     cust.Id,
			expErr: nil,
			expLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, mr.Delete(tc.id), tc.expErr)
			ll, _ := mr.ListByUser(cust.UserId)
			assert.Len(t, ll, tc.expLen)
		})
	}
}
