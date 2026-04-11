package postgres

import (
	"goddd/domain/models"
	"goddd/pkg/errors"
	"log"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func clearDB(tb testing.TB) {
	if _, err := db.Exec("DELETE FROM customer"); err != nil {
		tb.Fatal(err)
	}
	log.Println("teardown test")
}

func addCustomers(t *testing.T, user models.User, customers ...models.Customer) func(testing.TB) {
	// add customers
	tx := db.MustBegin()
	// add user
	_, err := tx.NamedExec("INSERT INTO \"user\" (id, username, password) VALUES (:id, :username, :password)", &user)
	require.NoError(t, err)
	_, err = tx.NamedExec("INSERT INTO customer (id, user_id, name, email) VALUES (:id, :user_id, :name, :email)", customers)
	require.NoError(t, err)
	if err := tx.Commit(); err != nil {
		t.Error(err)
	}

	return func(tb testing.TB) {
		//clear repo
		clearDB(tb)
	}
}

func TestCustomerRepo_GetCustomer(t *testing.T) {
	// fake user to add to repository
	u := models.User{Id: "u_ID", Username: "Jonathan", Password: "jonathan_password"}
	// Create a fake customer to add to repository
	cust := models.Customer{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}
	// Create the repo to use, and add some test Data to it for testing
	r := customerRepo{db}

	defer addCustomers(t, u, cust)(t)

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
			c, err := r.Get(tc.id)
			if err == nil {
				assert.Nil(t, tc.expErr)
				assert.Equal(t, c, cust)
			} else {
				assert.Equal(t, errors.GetID(err), errors.GetID(tc.expErr))
				assert.Equal(t, err.Error(), tc.expErr.Error())
			}
		})
	}
}

func TestCustomerRepo_AddCustomer(t *testing.T) {
	cust := models.Customer{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}
	r := customerRepo{db}
	defer clearDB(t)

	// no item yet
	var count int
	err := db.QueryRow("SELECT count(*) FROM customer").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, count, 0, "repository should be empty")

	// add a new customer
	assert.NoError(t, r.Add(cust), "unexpected error")

	err = db.QueryRow("SELECT count(*) FROM customer").Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, count, 1, "repository should contain a customer")

	// add the same customer twice should fail
	err = r.Add(cust)
	if err == nil {
		t.Fatal(err)
	}
	assert.Equal(t, errors.GetID(err), errors.RepoItemAlreadyExist)
	assert.Equal(t, err.Error(), "customer already exist in repository")
}

func TestCustomerRepo_ListCustomers(t *testing.T) {
	// fake user to add to repository
	u := models.User{Id: "u_ID", Username: "Jonathan", Password: "jonathan_password"}
	cust := models.Customer{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}
	r := customerRepo{db}

	// DB should be empty
	ll, err := r.ListByUser(cust.UserId)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, ll, 0, "repository should be empty")

	defer addCustomers(t, u, cust)(t)

	ll, err = r.ListByUser(cust.UserId)
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, ll, 1, "repository should have a single element")
}

func TestCustomerRepo_DeleteCustomers(t *testing.T) {
	// fake user to add to repository
	u := models.User{Id: "u_ID", Username: "Jonathan", Password: "jonathan_password"}
	cust := models.Customer{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}
	r := customerRepo{db}

	testCases := []struct {
		name   string
		id     string
		expErr error
		expLen int
	}{
		{
			name:   "empty customer id",
			id:     "",
			expErr: errors.InternalError.New("customer id must not be empty"),
			expLen: 1,
		},
		{
			name:   "unknow customer id",
			id:     "unknow_id",
			expErr: nil,
			expLen: 1,
		},
		{
			name:   "delete customer id Ok",
			id:     cust.Id,
			expErr: nil,
			expLen: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			defer addCustomers(t, u, cust)(t)
			err := r.Delete(tc.id)
			assert.Equal(t, err, tc.expErr)
			ll, err := r.ListByUser(cust.UserId)
			if err != nil {
				t.Fatal(err)
			}
			assert.Len(t, ll, tc.expLen, "repository should have a single element")
		})
	}
}
