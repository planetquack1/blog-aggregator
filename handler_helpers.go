package main

import (
	"encoding/json"
	"net/http"
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
