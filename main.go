package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/planetquack1/blog-aggregator/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}

func main() {

	// Load environment variables
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("PORT is not found in the environment")
	}
	dbURLString := os.Getenv("DB_URL")
	if dbURLString == "" {
		log.Fatal("DB_URL is not found in the environment")
	}

	// Open a connection and add it to config
	conn, err := sql.Open("postgres", dbURLString)
	if err != nil {
		log.Fatal("Can't connect to database")
	}

	cfg := apiConfig{
		DB: database.New(conn),
	}

	// Initialize ServeMux
	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerError)

	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.getUsers))
	mux.HandleFunc("POST /v1/users", cfg.postUsers)

	mux.HandleFunc("GET /v1/feeds", cfg.getFeeds)
	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.postFeeds))

	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.getFeedFollows))
	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.postFeedFollows))
	mux.HandleFunc("DELETE /v1/feed_follows/{feedFollowID}", cfg.middlewareAuth(cfg.deleteFeedFollows))

	srv := http.Server{
		Addr:    ":" + portString,
		Handler: mux,
	}

	srv.ListenAndServe()

}
