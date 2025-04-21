package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/arturogood17/Chirpy/internal/database"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	serverHits atomic.Int32
	Queries    *database.Queries
}

func main() {
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error open channel to db - %v", err)
	}
	dbQueries := database.New(db)
	const port = "8080"
	cfg := apiConfig{
		Queries: dbQueries,
	}
	mux := http.NewServeMux()
	srvr := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	mux.Handle("/app/", cfg.middlewareConfig(http.StripPrefix("/app", http.FileServer(http.Dir("."))))) //sirve todo lo que están en las carpetas
	mux.HandleFunc("GET /admin/metrics", cfg.handlerServerHits)                                         //se agrego el path admin para uso interno
	mux.HandleFunc("POST /admin/reset", cfg.handlerServerHitsReset)
	mux.HandleFunc("GET /api/healthz", hReadiness)
	mux.HandleFunc("POST /api/validate_chirp", hValChirpy)
	err = srvr.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't listen and serve from server - %v", err)
	}
}
