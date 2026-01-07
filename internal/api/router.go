package api

import (
	"net/http"

	"github.com/TheMaru/training-organiser/internal/auth"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	_ "github.com/TheMaru/training-organiser/docs"
	"github.com/swaggo/http-swagger/v2"
)

func NewRouter(cfg *ApiConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://next-app:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:8080/swagger/doc.json")))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	v1Router := chi.NewRouter()

	v1Router.Post("/users", cfg.HandleRegisterUser)
	v1Router.Post("/login", cfg.HandleLogin)
	v1Router.Post("/refresh", cfg.HandleRefreshToken)
	v1Router.Post("/revoke", cfg.HandleRevokeToken)

	v1Router.Group(func(r chi.Router) {
		r.Use(auth.MiddlewareAuth)

		r.Get("/users", cfg.HandleListUsers)
		r.Get("/users/{id}", cfg.HandleListUser)
		r.Get("/users/me", cfg.HandleGetMyProfile)
		r.Put("/users/{id}", cfg.HandleUpdateUser)
		r.Delete("/users/{id}", cfg.HandleDeleteUser)
	})

	r.Mount("/v1", v1Router)

	return r
}
