package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

func handlerReadiness(w http.ResponseWriter, r *http.Request) {

	status := Status{
		Status: "OK",
	}

	respondWithJSON(w, http.StatusOK, status)
}

func handlerError(w http.ResponseWriter, r *http.Request) {

	respondWithError(w, http.StatusInternalServerError, "Internal Server Error")
}

func (cfg *apiConfig) postUsers(w http.ResponseWriter, r *http.Request) {

	// Read in struct
	name, err := decodeJSON(r, Name{})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to parse request")
	}

	currentTime := time.Now()

	userParams := CreateUserParams{
		ID:        uuid.NewString(),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
		Name:      name,
	}

	cfg.DB.CreateUser(r.Context(), userParams)

	// Marshall and write
	dat, err := json.Marshal(userParams)
	if err != nil {
		log.Printf("Error marshalling user: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dat)
}
