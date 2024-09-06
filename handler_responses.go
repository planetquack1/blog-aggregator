package main

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {

	dat, err := json.Marshal(payload)
	if err != nil {
		w.Write([]byte("Failed to marshal JSON"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {

	error_msg := Error{
		Error: msg,
	}

	respondWithJSON(w, code, error_msg)
}
