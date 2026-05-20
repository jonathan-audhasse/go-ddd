package httapi

import (
	"goddd/api/http/handler"
	"goddd/internal/application"
	"net/http"

	appmdw "goddd/api/http/middleware"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(services *application.Services) http.Handler {
	r := chi.NewRouter()

	// middleware
	// r.Use(middleware.RequestID)
	// r.Use(middleware.RealIP)
	r.Use(appmdw.LoggerHandler)
	r.Use(middleware.Recoverer)

	// use CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Accept", "Authorization", "Content-Type"},
	}))

	healthHandler := handler.NewHealthHandler(services.Health)
	userHandler := handler.NewUserHandler(services.User)

	// routes
	r.Get("/health", healthHandler.IsHealthy)

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/{id}", userHandler.GetUser)
	})

	return r
}

// func NewRouter(services *application.Services) http.Handler {
// 	r := chi.NewRouter()

// 	// global middleware
// 	// r.Use(middleware.RequestID)
// 	// r.Use(middleware.RealIP)
// 	r.Use(middleware.Logger)
// 	r.Use(middleware.Recoverer)

// 	healthHandler := handler.NewHealthHandler(services.Health)
// 	userHandler := handler.NewUserHandler(services.User)

// 	// public routes
// 	r.Group(func(r chi.Router) {
// 		r.Post("/login", NewAuthHandler(svc.UserService, jwtSvc).Login)
// 	})
// 	// routes
// 	r.Get("/health", healthHandler.IsHealthy)

// 	jwtSvc :=
// 	// protected routes
// 	r.Group(func(r chi.Router) {
// 		r.Use(appmw.JWTAuth(jwtSvc))

// 		r.Route("/users", func(r chi.Router) {
// 			r.Post("/", userHandler.CreateUser)
// 			r.Get("/{id}", userHandler.GetUser)
// 		})
// 	})

// 	return r
// }
