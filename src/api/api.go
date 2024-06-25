package api

import (
	"context"
	"fmt"
	"goddd/src/infra/repositories/memory"
	"goddd/src/services"
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
	repo := memory.NewRepository()
	// init services
	services, err := services.NewServices(repo)
	if err != nil {
		log.Panic(err)
	}
	return Api{router: NewRouter(services), cfg: NewConfig()}
}

// server
func (a Api) ServeBis() {
	// server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", a.cfg.Port),
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

type Cust struct {
	Id    string
	Email string
	Name  string `db:"name"`
}

func (a Api) Serve() {
	db, err := sqlx.Connect("postgres", a.cfg.DbUrl())
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nUnable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()
	fmt.Printf("connected. db=%v\n", db)
	var data Cust
	if err := db.Get(&data, "select * from customer"); err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%T\n%v\n%#v", data, data, data)
}
