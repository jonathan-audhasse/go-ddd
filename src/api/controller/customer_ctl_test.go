package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goddd/src/domain/models"
	"goddd/src/infra/repositories/memory"
	"goddd/src/services"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// func TestMain(m *testing.M) {
//   // Setup code goes here
//   code := m.Run()
//   // Teardown code goes here
//   os.Exit(code)
// }

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

func addCustomers(repo *memory.CustomerMemRepo, customers []models.Customer) func(testing.TB) {
	// add customers
	for _, c := range customers {
		repo.Add(c)
	}

	return func(tb testing.TB) {
		// empty repo
		assert.NoError(tb, repo.Delete(), "no error expected")
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
			customers: []models.Customer{{ID: "an_ID", Name: "Johnny", Email: "a@b.c"}},
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
	defer addCustomers(repo, []models.Customer{{ID: "an_ID", Name: "Johnny", Email: "a@b.c"}})(t)
	testCases := []struct {
		name       string
		id         string
		statusCode int
		expRes     map[string]any
	}{
		{
			name:       "get customer (ID not found)",
			id:         "unknow_id",
			statusCode: http.StatusNotFound,
			expRes:     map[string]any{"errorId": float64(3000), "message": "(3000) fail to get customer (id=unknow_id): customer (id=unknow_id) not found in repository"},
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

// POST /service Test
// func TestWordApi_AddWord(t *testing.T) {
// 	repo := memory.NewWordMemoryRepository()
// 	service, err := services.NewWordService(repo)
// 	if err != nil {
// 		t.Errorf("unexpected error, got '%v'", err)
// 	}
// 	api := &WordApi{service}
// 	t.Run("tests POST /service", func(t *testing.T) {
// 		r := gin.Default()
// 		r.POST("/service", api.AddWord)

// 		testCases := []struct {
// 			name       string
// 			body       map[string]any
// 			statusCode int
// 			expRes     map[string]any
// 		}{
// 			{
// 				name:       "(KO) request body null (status code 400)",
// 				body:       nil,
// 				statusCode: 400,
// 				expRes:     map[string]any{"message": "body should be of format {'word': 'abc'}"},
// 			},
// 			{
// 				name:       "(KO) request body empty (status code 400)",
// 				body:       map[string]any(nil),
// 				statusCode: 400,
// 				expRes:     map[string]any{"message": "body should be of format {'word': 'abc'}"},
// 			},
// 			{
// 				name:       "(KO) request body wrong json format (status code 400)",
// 				body:       map[string]any{"w": 1},
// 				statusCode: 400,
// 				expRes:     map[string]any{"message": "body should be of format {'word': 'abc'}"},
// 			},
// 			{
// 				name:       "(KO) request body wrong word type (status code 400)",
// 				body:       map[string]any{"word": 2},
// 				statusCode: 400,
// 				expRes:     map[string]any{"message": "body should be of format {'word': 'abc'}"},
// 			},
// 			{
// 				name:       "(KO) request body invalid word (status code 400)",
// 				body:       map[string]any{"word": "a b"},
// 				statusCode: 422,
// 				expRes:     map[string]any{"message": "a word should contains only alphabetic characters and should not be empty"},
// 			},
// 			{
// 				name:       "(Success) should add a new word (status code 201)",
// 				body:       map[string]any{"word": "aBc"},
// 				statusCode: 201,
// 				expRes:     map[string]any{"frequency": float64(1), "word": "abc"},
// 			},
// 			{
// 				name:       "(Success) should update existing word (status code 201)",
// 				body:       map[string]any{"word": "ABC"},
// 				statusCode: 201,
// 				expRes:     map[string]any{"frequency": float64(2), "word": "abc"},
// 			},
// 		}
// 		// run all test cases
// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				b := new(bytes.Buffer)
// 				json.NewEncoder(b).Encode(tc.body)
// 				req, err := http.NewRequest(http.MethodPost, "/service", b)
// 				if err != nil {
// 					t.Errorf("unexpected error, got '%v'", err)
// 				}
// 				rec := httptest.NewRecorder()
// 				r.ServeHTTP(rec, req)

// 				// check status code
// 				assert.Equal(t, tc.statusCode, rec.Code)
// 				res := make(map[string]any)
// 				// check error message
// 				if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
// 					t.Errorf("unexpected error, got '%v'", err)
// 				}

// 				assert.Equal(t, tc.expRes, res)
// 			})
// 		}
// 	})
// }

// // GET /service Test
// func TestWordApi_GetWordFromPrefix(t *testing.T) {
// 	repo := memory.NewWordMemoryRepository()
// 	// populate repository
// 	repo.Add(models.Word{Value: "abc", Frequency: 2})
// 	service, err := services.NewWordService(repo)
// 	if err != nil {
// 		t.Errorf("unexpected error, got '%v'", err)
// 	}
// 	api := &WordApi{service}
// 	t.Run("tests GET /service", func(t *testing.T) {
// 		r := gin.Default()
// 		r.GET("/service", api.GetWordFromPrefix)

// 		testCases := []struct {
// 			name       string
// 			params     url.Values
// 			statusCode int
// 			expRes     any
// 		}{
// 			{
// 				name:       "(KO) missing 'prefix' query (status code 400)",
// 				params:     url.Values{},
// 				statusCode: 400,
// 				expRes:     map[string]string{"message": "missing 'prefix' parameter"},
// 			},
// 			{
// 				name:       "(KO) prefix query not valid (status code 400)",
// 				params:     url.Values{"prefix": {"a b"}},
// 				statusCode: 422,
// 				expRes:     map[string]string{"message": "a word should contains only alphabetic characters and should not be empty"},
// 			},
// 			{
// 				name:       "(Success) should return null when no word is found (status code 200)",
// 				params:     url.Values{"prefix": {"unknow"}},
// 				statusCode: 200,
// 				expRes:     map[string]string(nil),
// 			},
// 			{
// 				name:       "(Success) should return the word if found (status code 200)",
// 				params:     url.Values{"prefix": {"ab"}},
// 				statusCode: 200,
// 				expRes:     map[string]string{"word": "abc"},
// 			},
// 		}
// 		// run all test cases
// 		for _, tc := range testCases {
// 			t.Run(tc.name, func(t *testing.T) {
// 				req, err := http.NewRequest(http.MethodGet, "/service", nil)
// 				if err != nil {
// 					t.Errorf("unexpected error, got '%v'", err)
// 				}
// 				req.URL.RawQuery = tc.params.Encode()
// 				rec := httptest.NewRecorder()
// 				r.ServeHTTP(rec, req)

// 				// check status code
// 				assert.Equal(t, tc.statusCode, rec.Code)
// 				res := make(map[string]string)
// 				// check error message
// 				if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
// 					t.Errorf("unexpected error, got '%v'", err)
// 				}

// 				assert.Equal(t, tc.expRes, res)
// 			})
// 		}
// 	})
// }
