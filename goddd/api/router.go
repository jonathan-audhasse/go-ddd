package api

import (
	"goddd/api/controller"
	"goddd/internal/application"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func NewRouter(services *application.Services) http.Handler {
	r := chi.NewRouter()

	// middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	healthController := controller.NewHealthController(services.Health)
	userController := controller.NewUserController(services.User)

	// routes
	r.Get("/health", healthController.Health)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userController.CreateUser)
		r.Get("/{id}", userController.GetUser)
	})

	return r
}
