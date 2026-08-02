package httapi

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goddd/internal/application"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

// API to holds the API
type API struct {
	server *http.Server
	db     *pgxpool.Pool
}

// NewAPI creates a new API
func NewAPI(apiPort int, services *application.Services, db *pgxpool.Pool) *API {

	// router
	router := NewRouter(services)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", apiPort),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &API{
		server: srv,
		db:     db,
	}
}

// Close closes the API properly
func (a *API) Close(ctx context.Context) error {
	a.db.Close()
	return nil
}

// Serve serves the API
func (a *API) Serve() error {
	// run server in goroutine
	go func() {
		log.Info().Str("addr", a.server.Addr).Msg("starting http server")

		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("server crashed")
		}
	}()

	// wait for interrupt
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("shutting down server...")

	// timeout of 5s to let the server finished the current unfinished request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// shutdown server
	if err := a.server.Shutdown(ctx); err != nil {
		log.Err(err).Msg("failed to shutdown server cleanly")
	}

	// close db
	_ = a.Close(ctx)

	log.Info().Msg("server stopped")
	return nil
}
