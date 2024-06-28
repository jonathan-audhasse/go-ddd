package api

import (
	"context"
	"fmt"
	"goddd/src/domain/services"
	"goddd/src/infra/repositories/postgres"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Api struct {
	router *gin.Engine
	cfg    Config
}

func NewApi() Api {
	// init all repositories
	cfg := NewConfig()
	url := cfg.DbUrl()
	// repo := memory.NewRepository()
	log.Println("Set up DB. Connecting...")
	db, err := sqlx.Connect("postgres", url)
	if err != nil {
		log.Fatalf("fail to connect to DB. err=%s", err)
	}
	repo := postgres.NewRepository(db)
	// init services
	services, err := services.NewServices(repo)
	if err != nil {
		log.Panic(err)
	}
	return Api{router: NewRouter(services), cfg: cfg}
}

// server
func (a Api) Serve() {
	// server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.ApiPort),
		Handler: a.router.Handler(),
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		log.Printf("Listen server (port %s)...\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout
	// of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}
