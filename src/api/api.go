package api

import (
	"context"
	"fmt"
	"goddd/domain/services"
	"goddd/infra/repositories/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
)

type Api struct {
	server  *http.Server
	closeDb func()
}

func NewApi(cfg Config) Api {
	// repo := memory.NewRepository()
	repo, closeDb, err := postgres.NewRepository()
	if err != nil {
		log.Fatalf("failed to connect to DB. err=%s", err)
	}
	// init services
	services := services.NewServices(repo)
	if err != nil {
		log.Panic(err)
	}
	return Api{
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", cfg.ApiPort),
			Handler: NewRouter(services).Handler(),
		},
		closeDb: closeDb,
	}
}

// server
func (a Api) Serve() {
	// run server in a thread
	go func() {
		log.Printf("Listen server (port %s)...\n", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// set a timeout of 5 seconds to let the server handle the current unfinished request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// force server shutdown
	if err := a.server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	// close DB
	a.closeDb()
	log.Println("Server exiting")
}
