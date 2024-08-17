package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Erorr marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(code)
	w.Write(data)
}

func responseWithError(w http.ResponseWriter, code int, msg string) {
	// Create an error response payload
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResponse struct {
		Error string `json:"error"`
	}

	responseWithJSON(w, code, errorResponse{
		Error: msg,
	})
}
