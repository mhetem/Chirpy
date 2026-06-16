package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type jResponse struct {
	Body string `json:"body"`
}

type jError struct {
	Error string `json:"error"`
}

type jValid struct {
	Valid bool `json:"valid"`
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	payload := jError{
		Error: msg,
	}
	respondWithJSON(w, code, payload)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	respBody := payload
	dat, err := json.Marshal(respBody)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
