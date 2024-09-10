package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/planetquack1/blog-aggregator/internal/database"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {

	type Status struct {
		Status string `json:"status"`
	}

	status := Status{
		Status: "OK",
	}

	respondWithJSON(w, http.StatusOK, status)
}

func handlerError(w http.ResponseWriter, r *http.Request) {

	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (cfg *apiConfig) postUsers(w http.ResponseWriter, r *http.Request) {

	type Name struct {
		Name string `json:"name"`
	}

	// Read in struct
	name, err := decodeJSON(r, Name{})
	if err != nil {
		respondWithError(w, 400, "Failed to parse JSON")
		return
	}

	currentTime := time.Now()

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate UUID")
		return
	}

	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        randomUUID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      name.Name,
	})
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Marshall and write
	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}

func (cfg *apiConfig) getUsers(w http.ResponseWriter, r *http.Request) {

	// Extract API key from the Authorization header
	apiKey := getAuthFromHeader(r, "ApiKey")

	user, err := cfg.DB.GetUser(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to get user")
		return
	}

	// Marshall and write
	dat, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
