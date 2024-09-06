package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/planetquack1/blog-aggregator.git/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	// Load environment variables
	godotenv.Load()
	port := os.Getenv("PORT")

	// Open a connection and add it to config
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal()
	}

	dbQueries := database.New(db)
	cfg := apiConfig{
		DB: dbQueries,
	}

	// Initialize ServeMux
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerError)
	mux.HandleFunc("POST /v1/users", cfg.postUsers)

	srv := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	srv.ListenAndServe()

}
