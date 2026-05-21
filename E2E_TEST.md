
For your architecture, the cleanest E2E strategy is:

Start the real app (router + services + DB) against a real Postgres Testcontainer, then hit HTTP endpoints with httptest.

This gives you:

✔ real HTTP
✔ real DB
✔ real migrations
✔ real repositories/services
✔ no mocks
✔ fast enough
✔ production-like behavior
🧱 Recommended E2E structure
tests/
├── e2e/
│   ├── main_test.go
│   ├── helpers.go
│   ├── user_test.go
│   └── health_test.go
✅ Goal

Test full flow:

HTTP -> router -> controller -> service -> repo -> postgres
🚀 1. Expose your router cleanly

You should already have something like:

api/http/router/router.go
func NewRouter(ctrls *controller.Controllers) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Cors)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", ctrls.UserController.CreateUser)
		r.Get("/{id}", ctrls.UserController.GetUser)
	})

	r.Get("/health", ctrls.HealthController.Health)

	return r
}
🚀 2. Build app bootstrap for tests

You already have good separation (bootstrap/database.go, bootstrap/services.go).

Now expose a test app builder.

✅ Example E2E bootstrap
tests/e2e/main_test.go
package e2e_test

import (
	"context"
	"net/http"
	"os"
	"testing"

	"goddd/internal/bootstrap"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	testServer *http.Server
	testDB     *pgxpool.Pool
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	app, cleanup, err := bootstrap.NewTestApplication(ctx)
	if err != nil {
		panic(err)
	}

	testServer = app.Server
	testDB = app.DB

	code := m.Run()

	cleanup()

	os.Exit(code)
}
🚀 3. Create bootstrap test app
internal/bootstrap/testing.go
type TestApplication struct {
	Server *http.Server
	DB     *pgxpool.Pool
}

func NewTestApplication(ctx context.Context) (*TestApplication, func(), error) {
	pg, err := newPGContainer(ctx)
	if err != nil {
		return nil, nil, err
	}

	db := pg.pool

	repos := repository.New(db)

	tm := transaction.NewTransactionManager(db)

	services := application.NewServices(repos, tm)

	ctrls := controller.NewControllers(services)

	router := router.NewRouter(ctrls)

	server := &http.Server{
		Handler: router,
	}

	cleanup := func() {
		db.Close()
		pg.Close(ctx)
	}

	return &TestApplication{
		Server: server,
		DB:     db,
	}, cleanup, nil
}
🚀 4. Use httptest
tests/e2e/user_test.go
package e2e_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateUser(t *testing.T) {
	body := map[string]any{
		"username": "john",
		"email":    "john@test.com",
	}

	b, _ := json.Marshal(body)

	req := httptest.NewRequest(
		http.MethodPost,
		"/users",
		bytes.NewReader(b),
	)

	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	testServer.Handler.ServeHTTP(rec, req)

	require.Equal(t, http.StatusCreated, rec.Code)
}
🔥 BEST PRACTICE: transaction isolation per test

Very important.

For E2E:

Option A (simple)
recreate DB/container once per suite
Option B (better)
rollback transaction per test

But for HTTP E2E this is harder because request lifecycle owns DB access.

So usually:

✅ Recommended for your architecture
Use truncate cleanup

Create helper:

func cleanupTables(t *testing.T) {
	_, err := testDB.Exec(context.Background(), `
		TRUNCATE TABLE users RESTART IDENTITY CASCADE
	`)
	require.NoError(t, err)
}

Then:

func TestCreateUser(t *testing.T) {
	cleanupTables(t)
🚀 BEST setup overall
Test Type	Uses
unit	mocks
integration	real repo + testcontainer
e2e	full HTTP + DB

You now have all 3 levels.

🧠 My recommendation for YOUR project
Keep:
Unit tests
services/controllers with mocks
Integration tests
repositories + postgres testcontainer
E2E tests
full HTTP flow

That’s an excellent architecture.

🎯 Final recommendation

Create:

tests/
  integration/
  e2e/

And keep infra tests out of internal/.

This scales MUCH better long-term.

--------------------


Yes — for your project, that is actually a very good approach.

Given your architecture and goals, the best setup is usually:

One real PostgreSQL instance + dedicated schema/database for tests

instead of spinning up Testcontainers for every E2E run.

This is very common in mature backend teams because it gives:

⚡ much faster test runs
🧪 stable environment
🔁 easier local iteration
🐳 simpler Docker Compose workflow
🚀 better CI performance
🧠 Recommended strategy for YOUR project

You already have:

migrations
schema support
clean infra separation
Docker Compose

So you are in a perfect position to use:

✅ shared Postgres instance

with:

public          -> dev
goddd_test      -> e2e

OR separate DBs:

goddd_dev
goddd_test
🎯 My recommendation
Prefer:
✅ Separate DATABASE (best)

instead of schema-only.

Why?

Because:

cleaner isolation
easier reset
no schema leakage
safer migrations
simpler tooling
mirrors production better
🧱 Recommended Docker setup
docker-compose.yml
services:
  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: abc123
      POSTGRES_DB: goddd_dev
    ports:
      - "5432:5432"

Then manually create:

CREATE DATABASE goddd_test;
🧪 E2E config
.env.test
DATABASE_URL=postgres://admin:abc123@localhost:5432/goddd_test?sslmode=disable
🚀 E2E flow

Before tests:

go run cmd/migrate/main.go -env=test -cmd=up

Run tests:

go test ./tests/e2e/...

After tests:

TRUNCATE users CASCADE;

or rollback transaction.

🔥 Even better: transaction-per-test

This is the ideal setup.

Each E2E test:

tx, _ := db.Begin(ctx)
defer tx.Rollback(ctx)

But for HTTP E2E you must inject tx into request context or use a test transaction middleware.

That becomes more advanced.

✅ Simpler and very effective approach

For now:

Use table cleanup between tests

Example:

func cleanupDB(t *testing.T) {
	_, err := testDB.Exec(context.Background(), `
		TRUNCATE TABLE users RESTART IDENTITY CASCADE
	`)
	require.NoError(t, err)
}

Fast and reliable.

⚖️ Testcontainers vs Real Test DB
	Testcontainers	Dedicated Test DB
isolation	excellent	good
startup speed	slower	fast
CI portability	excellent	requires infra
local DX	medium	excellent
production realism	high	high
maintenance	medium	low
🧠 What mature teams usually do
Local dev

👉 shared *_test database

CI

👉 ephemeral containerized DB

Best of both worlds.

🔥 My recommendation specifically for YOU

Given your architecture:

✅ Use:
Docker Compose Postgres
dedicated goddd_test database
migration runner
cleanup helpers
❌ Avoid:
starting Testcontainers for all E2E tests

Keep Testcontainers only for:

repository integration tests
CI isolation if needed
🧱 Recommended final structure
tests/
├── integration/     # repo tests with testcontainers
├── e2e/             # real app + real test DB
└── fixtures/
🎯 Final answer

Yes — for your project, using:

a real PostgreSQL instance with a dedicated test database

is likely the best E2E strategy.

It’s faster, simpler, and scales better than full Testcontainers E2E.

If you want, I can next help you build:

🔥 a complete E2E setup:
.env.test
Makefile
docker-compose
migration automation
DB cleanup helpers
GitHub Actions CI

