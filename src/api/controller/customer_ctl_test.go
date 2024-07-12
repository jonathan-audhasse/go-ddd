package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goddd/src/domain/models"
	"goddd/src/domain/services"
	"goddd/src/infra/repositories/postgres"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	// establish connection
	conn, close, err := postgres.Connect()

	if err != nil {
		log.Fatalf("failed to connect to DB. err=%s", err)
	}
	db = conn
	code := m.Run()
	// Teardown code goes here
	// close db session
	close()
	os.Exit(code)
}

// You can use *testing.T, if you want to test the code without benchmarking
func setupTest(tb testing.TB, method, relativePath, url string, handler gin.HandlerFunc, body io.Reader) *httptest.ResponseRecorder {
	log.Println("setup test")

	// set up engine
	r := gin.Default()
	// r.GET(url, handler) == r.Handle("GET", url, handler)
	r.Handle(method, relativePath, handler)

	// set up the request recorder
	// http.POST(url, body) == htt.NewRequest("POST", url, body)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		tb.Fatal(err)
	}
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)

	// return the response recorder
	return rec
}

func clearUsers(tb testing.TB) {
	if _, err := db.Exec("DELETE FROM \"user\""); err != nil {
		tb.Fatal(err)
	}
	fmt.Println("teardown test")
}

func addUsers(t *testing.T, users ...models.User) func(testing.TB) {
	// add users
	tx := db.MustBegin()
	tx.NamedExec("INSERT INTO \"user\" (id, username, password) VALUES (:id, :username, :password)", users)
	if err := tx.Commit(); err != nil {
		t.Error(err)
	}

	return func(tb testing.TB) {
		//clear repo
		clearUsers(tb)
	}
}

func clearCustomers(tb testing.TB) {
	if _, err := db.Exec("DELETE FROM customer"); err != nil {
		tb.Fatal(err)
	}
	fmt.Println("teardown test")
}

func addCustomers(t *testing.T, customers []models.Customer) func(testing.TB) {
	tx := db.MustBegin()
	// add customers
	tx.NamedExec("INSERT INTO customer (id, user_id, name, email) VALUES (:id, :user_id, :name, :email)", customers)
	if err := tx.Commit(); err != nil {
		t.Error(err)
	}

	return func(tb testing.TB) {
		// empty repo
		clearCustomers(tb)
	}
}

func TestCustomerCtl_ListCustomer(t *testing.T) {
	// Create a fake user to add to repository
	u := models.User{Id: "u_ID", Username: "Jonathan", Password: "jonathan_password"}
	defer addUsers(t, u)(t)

	// Create repo Repo
	repo := postgres.NewCustomerRepo(db)
	service := services.NewCustomerService(repo)
	testCases := []struct {
		customers []models.Customer
		name      string
		expRes    []map[string]string
	}{
		{
			customers: []models.Customer{},
			name:      "list customers (empty repository)",
			expRes:    []map[string]string{},
		},
		{
			customers: []models.Customer{{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}},
			name:      "list customers with a single customer",
			expRes:    []map[string]string{{"id": "an_ID", "name": "Johnny", "email": "a@b.c"}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// populate repo
			defer addCustomers(t, tc.customers)(t)

			rec := setupTest(t, http.MethodGet, "/customers", "/customers", (&CustomerCtl{service}).ListCustomers, nil)

			// check status code
			assert.Equal(t, http.StatusOK, rec.Code)
			res := make([]map[string]string, 0)

			if assert.NoError(t, json.NewDecoder(rec.Body).Decode(&res), "unexpected error") {
				assert.Equal(t, tc.expRes, res)
			}
		})
	}
}

func TestCustomerCtl_GetCustomer(t *testing.T) {
	// Create a fake user to add to repository
	u := models.User{Id: "u_ID", Username: "Jonathan", Password: "jonathan_password"}
	defer addUsers(t, u)(t)

	// Create repo Repo
	repo := postgres.NewCustomerRepo(db)
	service := services.NewCustomerService(repo)

	defer addCustomers(t, []models.Customer{{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}})(t)
	testCases := []struct {
		name       string
		id         string
		statusCode int
		expRes     map[string]any
	}{
		{
			name:       "get unknow customer",
			id:         "unknow_id",
			statusCode: http.StatusNotFound,
			expRes:     map[string]any{"errorId": float64(3000), "message": "(3000) fail to get customer id=unknow_id: customer id=unknow_id not found in repository"},
		},
		{
			name:       "get customer (ID ok)",
			id:         "an_ID",
			statusCode: http.StatusOK,
			expRes:     map[string]any{"id": "an_ID", "name": "Johnny", "email": "a@b.c"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := setupTest(t, http.MethodGet, "/customers/:id", fmt.Sprintf("/customers/%v", tc.id), (&CustomerCtl{service}).GetCustomer, nil)

			// check status code
			assert.Equal(t, tc.statusCode, rec.Code)
			res := make(map[string]any)

			if assert.NoError(t, json.NewDecoder(rec.Body).Decode(&res), "unexpected error") {
				assert.Equal(t, tc.expRes, res)
			}
		})
	}
}

func TestCustomerCtl_DeleteCustomer(t *testing.T) {
	// Create a fake user to add to repository
	u := models.User{Id: "u_ID", Username: "Jonathan", Password: "jonathan_password"}
	defer addUsers(t, u)(t)

	// Create repo Repo
	repo := postgres.NewCustomerRepo(db)
	service := services.NewCustomerService(repo)

	defer addCustomers(t, []models.Customer{{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}})(t)
	testCases := []struct {
		name       string
		id         string
		statusCode int
	}{
		{
			name:       "remove unknow customer",
			id:         "unknow_id",
			statusCode: http.StatusNoContent,
		},
		{
			name:       "remove customer (ID ok)",
			id:         "an_ID",
			statusCode: http.StatusNoContent,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := setupTest(t, http.MethodDelete, "/customers/:id", fmt.Sprintf("/customers/%v", tc.id), (&CustomerCtl{service}).DeleteCustomer, nil)

			// check status code
			assert.Equal(t, tc.statusCode, rec.Code)
		})
	}
}

func TestCustomerCtl_AddCustomer(t *testing.T) {
	// Create a fake user to add to repository
	u := models.User{Id: "u_ID", Username: "Jonathan", Password: "jonathan_password"}
	defer addUsers(t, u)(t)

	// Create repo Repo
	repo := postgres.NewCustomerRepo(db)
	service := services.NewCustomerService(repo)

	testCases := []struct {
		name       string
		body       map[string]string
		statusCode int
		expRes     map[string]any
	}{
		{
			name:       "add customer (empty body)",
			body:       map[string]string{},
			statusCode: http.StatusBadRequest,
			expRes:     map[string]any{"errorId": float64(2000), "message": "customer body should contain at least a name"},
		},
		{
			name:       "add customer OK",
			body:       map[string]string{"name": "Johnny", "email": "a@b.c"},
			statusCode: http.StatusCreated,
			expRes:     map[string]any{"name": "Johnny", "email": "a@b.c"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			jsonCust, _ := json.Marshal(tc.body)
			// var b bytes.Buffer // use &b as *io.Reader
			// json.NewEncoder(b).Encode(tc.body)
			rec := setupTest(t, http.MethodPost, "/customers", "/customers", (&CustomerCtl{service}).AddCustomer, bytes.NewBuffer(jsonCust))

			// check status code
			assert.Equal(t, tc.statusCode, rec.Code)
			res := make(map[string]any)

			if assert.NoError(t, json.NewDecoder(rec.Body).Decode(&res), "unexpected error") {
				switch tc.statusCode {
				case http.StatusCreated:
					assert.Contains(t, res, "id")
					assert.Equal(t, tc.expRes["name"], res["name"])
					assert.Equal(t, tc.expRes["email"], res["email"])
					ll, _ := repo.ListByUser(u.Id)
					assert.Len(t, ll, 1)
				default:
					// error
					assert.Equal(t, tc.expRes, res)
				}
			}
		})
	}
}
