package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync/atomic"

	"github.com/mhetem/Chirpy/internal/database"

	_ "github.com/lib/pq"

	"github.com/joho/godotenv"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handlerMetrics(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(fmt.Sprintf("<html>\n   <body>\n    <h1>Welcome, Chirpy Admin</h1>\n    <p>Chirpy has been visited %d times!</p>\n  </body>\n</html>", count)))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("Hits reset"))
}

func (cfg *apiConfig) handlerValidate(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := jResponse{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Error decoding")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
		return
	}
	cleaned_body := removeProfane(params.Body)
	type returnVals struct {
		CleanedBody string `json:"cleaned_body"`
	}
	resp := returnVals{
		CleanedBody: cleaned_body,
	}
	respondWithJSON(w, 200, resp)

}

func removeProfane(body string) string {
	words := strings.Split((body), " ")
	for i, word := range words {
		if strings.ToLower(word) == "kerfuffle" {
			words[i] = "****"
		}
		if strings.ToLower(word) == "sharbert" {
			words[i] = "****"
		}
		if strings.ToLower(word) == "fornax" {
			words[i] = "****"
		}
	}
	return strings.Join(words, " ")

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

	cfg := apiConfig{
		dbQueries: dbQueries,
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
	mux.HandleFunc("POST /api/validate_chirp", cfg.handlerValidate)

	serv := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = serv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}
