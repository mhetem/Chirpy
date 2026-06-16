package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/mhetem/Chirpy/internal/database"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	Platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load the DB")
	}
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Can't open the db")
	}

	dbQueries := database.New(db)
	platform := os.Getenv("PLATFORM")

	cfg := apiConfig{
		dbQueries: dbQueries,
		Platform:  platform,
	}
	mux := http.NewServeMux()

	handler := http.StripPrefix("/app", http.FileServer(http.Dir(".")))

	mux.Handle("/app/", cfg.middlewareMetricsInc(handler))

	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK\n"))
	})

	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/chirps", cfg.handlerChirps)
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("GET /api/chirps", cfg.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetSingleChirp)
	mux.HandleFunc("POST /api/login", cfg.handlerLogin)

	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = serv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
