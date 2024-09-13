package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/planetquack1/blog-aggregator/internal/database"
)

// decodeJSON takes in an http.Request and a struct type, decodes the JSON body into the provided struct type, and returns it.
func decodeJSON[T any](r *http.Request, target T) (T, error) {
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&target)
	if err != nil {
		return target, err
	}
	return target, nil
}

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract and validate API key from headers
		authHeader := r.Header.Get("Authorization")
		// Get user by API key
		user, err := cfg.DB.GetUser(r.Context(), authHeader)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Unauthorized")
			return
		}
		// If API key is valid, call the handler with the user data
		fmt.Println(user)
		handler(w, r, user)
	}
}
