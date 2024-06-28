package postgres

import (
	"context"
	"fmt"
	"goddd/src/domain/models"
	"goddd/src/pkg/errors"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/sethvargo/go-envconfig"
	"github.com/stretchr/testify/assert"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	// establish connection
	var cfg struct {
		Host string `env:"DB_HOST, default=db"`
		Port int    `env:"DB_PORT, default=5432"`
		Name string `env:"DB_NAME, default=postgres"`
		User string `env:"DB_USER, default=admin"`
		Pwd  string `env:"DB_PWD, default=abc123"`
	}
	if err := envconfig.Process(context.Background(), &cfg); err != nil {
		log.Fatal(err)
	}
	url := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Pwd, cfg.Name)
	log.Printf("Set up DB. Connecting to %s...\n", url)
	conn, err := sqlx.Connect("postgres", url)

	if err != nil {
		log.Fatalf("fail to connect to DB. err=%s", err)
	}
	db = conn
	code := m.Run()
	// Teardown code goes here
	if err := conn.Close(); err != nil {
		log.Panicf("Fail to close DB session. err=%v", err)
	}
	log.Printf("DB session Closed")
	os.Exit(code)
}

func emptyDB(tb testing.TB) {
	if _, err := db.Exec("DELETE FROM customer"); err != nil {
		tb.Fatal(err)
	}
	log.Println("teardown test")
}

func addCustomers(t *testing.T, customers ...models.Customer) func(testing.TB) {
	// add customers
	tx := db.MustBegin()
	tx.NamedExec("INSERT INTO customer (id, name, email) VALUES (:id, :name, :email)", customers)
	if err := tx.Commit(); err != nil {
		t.Error(err)
	}

	return func(tb testing.TB) {
		//empty repo
		emptyDB(tb)
	}
}

func TestCustomerPgRepo_GetCustomer(t *testing.T) {
	// Create a fake customer to add to repository
	cust := models.Customer{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}
	// Create the repo to use, and add some test Data to it for testing
	r := customerPgRepo{db}

	defer addCustomers(t, cust)(t)

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

func TestCustomerPgRepo_AddCustomer(t *testing.T) {
	cust := models.Customer{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}
	r := customerPgRepo{db}
	defer emptyDB(t)

	// no item yet
	var count int
	db.QueryRow("SELECT count(*) FROM customer").Scan(&count)
	assert.Equal(t, count, 0, "repository should be empty")

	// add a new customer
	assert.NoError(t, r.Add(cust), "unexpected error")

	db.QueryRow("SELECT count(*) FROM customer").Scan(&count)
	assert.Equal(t, count, 1, "repository should contain a customer")

	// add the same customer twice should fail
	err := r.Add(cust)
	if err == nil {
		t.Fatal(err)
	}
	assert.Equal(t, errors.GetID(err), errors.RepoItemAlreadyExist)
	assert.Equal(t, err.Error(), "customer already exist in repository")
}

func TestCustomerPgRepo_ListCustomers(t *testing.T) {
	cust := models.Customer{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}
	r := customerPgRepo{db}

	// DB should be empty
	ll, err := r.List()
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, ll, 0, "repository should be empty")

	defer addCustomers(t, cust)(t)

	ll, err = r.List()
	if err != nil {
		t.Fatal(err)
	}
	assert.Len(t, ll, 1, "repository should have a single element")
}

func TestCustomerPgRepo_DeleteCustomers(t *testing.T) {
	cust := models.Customer{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}
	r := customerPgRepo{db}

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
			defer addCustomers(t, cust)(t)
			err := r.Delete(tc.id)
			assert.Equal(t, err, tc.expErr)
			ll, err := r.List()
			if err != nil {
				t.Fatal(err)
			}
			assert.Len(t, ll, tc.expLen, "repository should have a single element")
		})
	}
}
