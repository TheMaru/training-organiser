package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TheMaru/training-organiser/internal/api"
	"github.com/TheMaru/training-organiser/internal/auth"
	"github.com/TheMaru/training-organiser/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func main() {
	godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	defer conn.Close()

	queries := database.New(conn)

	apiCfg := &api.ApiConfig{
		DB: queries,
	}

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

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})

	v1Router := chi.NewRouter()

	v1Router.Post("/users", apiCfg.HandleRegisterUser)
	v1Router.Post("/login", apiCfg.HandleLogin)
	v1Router.Post("/refresh", apiCfg.HandleRefreshToken)
	v1Router.Post("/revoke", apiCfg.HandleRevokeToken)

	r.Group(func(r chi.Router) {
		r.Use(auth.MiddlewareAuth)

		// auth routes here
	})

	r.Mount("/v1", v1Router)

	srv := &http.Server{
		Addr:         ":8080", // TODO: should this be configurable via env?
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("listening :8080")
	log.Fatal(srv.ListenAndServe())
}
