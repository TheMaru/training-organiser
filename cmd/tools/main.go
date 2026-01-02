package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/TheMaru/training-organiser/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	godotenv.Load()

	if len(os.Args) < 2 {
		fmt.Println("expected 'grant_admin' subcommand")
		os.Exit(1)
	}

	dbURL := os.Getenv("DATABASE_URL")
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't connect to database:", err)
	}
	defer conn.Close()

	queries := database.New(conn)

	switch os.Args[1] {
	case "grant_admin":
		grantAdmin(queries, os.Args[2])
	default:
		fmt.Println("command not existing")
		os.Exit(1)
	}
}

func grantAdmin(db *database.Queries, username string) {
	user, err := db.GetUserByUserName(context.Background(), username)
	if err != nil {
		log.Fatalf("User %s not found: %v", username, err)
	}

	err = db.GrantAdminRole(context.Background(), user.ID)
	if err != nil {
		log.Fatal("Couldn't grant admin role", err)
	}

	fmt.Println("Admint role granted")
}
