package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/TheMaru/training-organiser/internal/api"
	"github.com/TheMaru/training-organiser/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// @title Training Organiser API
// @version 1.0
// @description This is the API for the Training Organiser application
// @host localhost:8080
// @BasePath /v1
// @securityDefinitions.apiKey BearerAuth
// @in header
// @name Authorization
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

	r := api.NewRouter(apiCfg)

	srv := &http.Server{
		Addr:         ":8080", // TODO: should this be configurable via env?
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Println("listening :8080")
	log.Fatal(srv.ListenAndServe())
}
