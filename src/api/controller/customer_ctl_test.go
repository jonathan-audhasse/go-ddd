package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goddd/src/domain/models"
	"goddd/src/domain/repository"
	"goddd/src/domain/services"
	"goddd/src/infra/repositories/memory"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// func setupSuite(*testing.T) func(*testing.T) {
// 	log.Println("setup suite")

// 	return func(*testing.T) {
// 		log.Println("teardown suite")
// 	}
// }

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

func addCustomers(repo repository.CustomerRepository, customers []models.Customer) func(testing.TB) {
	// add customers
	for _, c := range customers {
		repo.Add(c)
	}

	return func(tb testing.TB) {
		// empty repo
		for _, c := range customers {
			assert.NoError(tb, repo.Delete(c.Id), "no error expected")
		}

		log.Println("teardown test")
	}
}

func TestCustomerCtl_ListCustomer(t *testing.T) {
	repo := memory.NewCustomerMemRepo()
	service, err := services.NewCustomerService(repo)
	if err != nil {
		t.Fatal(err)
	}
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
			defer addCustomers(repo, tc.customers)(t)

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
	repo := memory.NewCustomerMemRepo()
	service, err := services.NewCustomerService(repo)
	if err != nil {
		t.Fatal(err)
	}
	defer addCustomers(repo, []models.Customer{{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}})(t)
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
	repo := memory.NewCustomerMemRepo()
	service, err := services.NewCustomerService(repo)
	if err != nil {
		t.Fatal(err)
	}
	defer addCustomers(repo, []models.Customer{{Id: "an_ID", Name: "Johnny", Email: "a@b.c"}})(t)
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
	repo := memory.NewCustomerMemRepo()
	service, err := services.NewCustomerService(repo)
	if err != nil {
		t.Fatal(err)
	}
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
					ll, _ := repo.List()
					assert.Len(t, ll, 1)
				default:
					// error
					assert.Equal(t, tc.expRes, res)
				}
			}
		})
	}
}
