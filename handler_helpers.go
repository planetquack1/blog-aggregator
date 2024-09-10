package main

import (
	"encoding/json"
	"net/http"
	"strings"
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

func getAuthFromHeader(r *http.Request, prefix string) string {
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, (prefix + " ")) {
		// Return the token without the "Bearer " prefix
		return strings.TrimPrefix(authHeader, (prefix + " "))
	}
	// Return the header as is if "Bearer " is not found
	return authHeader
}
