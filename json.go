package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type jResponse struct {
	Body string `json:"body"`
}

type jError struct {
	Error string `json:"error"`
}

type paramChirp struct {
	Body   string    `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

type jChirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
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
