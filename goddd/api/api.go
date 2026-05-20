package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"goddd/api/controller"
	"goddd/internal/application"
	"goddd/internal/infrastructure/http/router"
	"goddd/internal/infrastructure/persistence/postgres"
	"goddd/internal/infrastructure/persistence/transaction"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type API struct {
	server *http.Server
	db     *pgxpool.Pool
}

func NewAPI(cfg Config) (*API, error) {
	ctx := context.Background()

	// db
	db, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// repositories
	repos := postgres.NewRepository(db)

	// transaction manager
	tm := transaction.NewTransactionManager(db)

	// services
	services := application.NewServices(repos, tm)

	// controllers
	healthCtrl := controller.NewHealthController(services.Health)
	userCtrl := controller.NewUserController(services.User)

	// router
	r := router.NewRouter(healthCtrl, userCtrl)

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.APIPort),
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &API{
		server: srv,
		db:     db,
	}, nil
}

func (a *API) Close(ctx context.Context) error {
	a.db.Close()
	return nil
}

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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown server cleanly")
	}

	_ = a.Close(ctx)

	log.Info().Msg("server stopped")
	return nil
}
