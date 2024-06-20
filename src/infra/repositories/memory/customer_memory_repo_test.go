package memory

import (
	"goddd/src/domain/models"
	"goddd/src/pkg/errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCustomerMemoryRepo_GetCustomer(t *testing.T) {
	// Create a fake customer to add to repository
	cust := models.Customer{ID: "an_ID", Name: "Johnny", Email: "a@b.c"}
	// Create the repo to use, and add some test Data to it for testing
	// Skip Factory for this
	mr := CustomerMemoryRepository{
		customers: map[string]models.Customer{cust.ID: cust},
	}

	testCases := []struct {
		name   string
		id     string
		expErr error
	}{
		{
			name:   "Customer not found",
			id:     "another_ID",
			expErr: errors.RepoItemNotFound.New("customer (id=another_ID) not found in repository"),
		},
		{
			name:   "Customer Ok",
			id:     cust.ID,
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

func TestCustomerMemoryRepo_AddCustomer(t *testing.T) {
	mr := CustomerMemoryRepository{customers: map[string]models.Customer{}}

	cust := models.Customer{ID: "anID", Name: "Johnny"}

	err := mr.Add(cust)
	if err != nil {
		t.Errorf("error not expected, got %v", err)
	}

	found, err := mr.Get(cust.ID)
	if err != nil {
		t.Fatal(err)
	}
	if found.ID != cust.ID {
		t.Errorf("expected %v, got %v", cust.ID, found.ID)
	}
}

func TestCustomerMemoryRepo_ListCustomers(t *testing.T) {
	mr := CustomerMemoryRepository{customers: map[string]models.Customer{}}
	cust := models.Customer{ID: "anID", Name: "Johnny"}

	// Test memory is empty
	ll, err := mr.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(ll) != 0 {
		t.Errorf("repository should be empty")
	}

	err = mr.Add(cust)
	if err != nil {
		t.Errorf("error not expected, got %v", err)
	}

	ll, err = mr.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(ll) != 1 {
		t.Errorf("repository should have a single element")
	}
}
