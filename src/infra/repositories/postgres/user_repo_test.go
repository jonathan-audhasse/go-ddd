package postgres

import (
	"fmt"
	"goddd/src/domain/models"
	"goddd/src/pkg/errors"
	"log"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var db *sqlx.DB

func TestMain(m *testing.M) {
	// establish connection
	conn, close, err := Connect()

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

func TestUserRepo_GetById(t *testing.T) {
	// Create a fake user to add to repository
	u := models.User{Id: "an_ID", Username: "Jonathan", Password: "jonathan_password"}
	// Create repo Repo
	r := userRepo{db}

	defer addUsers(t, u)(t)

	testCases := []struct {
		name   string
		id     string
		expErr error
	}{
		{
			name:   "User not found",
			id:     "unknow_ID",
			expErr: errors.RepoItemNotFound.New("user not found in repository"),
		},
		{
			name:   "User Ok",
			id:     u.Id,
			expErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := r.GetById(tc.id)
			assert.Equal(t, err, tc.expErr)
			if err == nil {
				assert.Equal(t, c, u)
			}
		})
	}
}

func TestUserRepo_GetByUsername(t *testing.T) {
	// Create a fake user to add to repository
	u := models.User{Id: "an_ID", Username: "Jonathan", Password: "jonathan_password"}
	// Create repo Repo
	r := userRepo{db}

	defer addUsers(t, u)(t)

	testCases := []struct {
		name     string
		username string
		expErr   error
	}{
		{
			name:     "User not found",
			username: "unknow_username",
			expErr:   errors.RepoItemNotFound.New("user not found in repository"),
		},
		{
			name:     "User Ok",
			username: u.Username,
			expErr:   nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c, err := r.GetByUsername(tc.username)
			assert.Equal(t, err, tc.expErr)
			if err == nil {
				assert.Equal(t, c, u)
			}
		})
	}
}
