package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/arturogood17/Chirpy/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	serverHits atomic.Int32
	Queries    *database.Queries
	PLATFORM   string
	SECRET     string
	POLKAKEY   string
}

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error open channel to db - %v", err)
	}
	dbQueries := database.New(db)
	const port = "8080"
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM must be set")
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("SECRET environment variable must be set")
	}
	polkakey := os.Getenv("POLKA_KEY")
	if secret == "" {
		log.Fatal("POLKA API KEY must be set")
	}
	cfg := apiConfig{
		Queries:  dbQueries,
		PLATFORM: platform,
		SECRET:   secret,
		POLKAKEY: polkakey,
	}
	mux := http.NewServeMux()
	srvr := http.Server{
		Handler: mux,
		Addr:    ":" + port,
	}
	mux.Handle("/app/", cfg.middlewareConfig(http.StripPrefix("/app", http.FileServer(http.Dir("."))))) //sirve todo lo que están en las carpetas
	mux.HandleFunc("GET /admin/metrics", cfg.handlerServerHits)                                         //se agrego el path admin para uso interno
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.hSingleChirp)
	mux.HandleFunc("GET /api/chirps", cfg.hListChirps)
	mux.HandleFunc("GET /api/healthz", hReadiness)
	mux.HandleFunc("POST /api/users", cfg.hUser)
	mux.HandleFunc("POST /api/chirps", cfg.hChirp)
	mux.HandleFunc("POST /admin/reset", cfg.handlerServerHitsReset)
	mux.HandleFunc("POST /api/login", cfg.hLogin)
	mux.HandleFunc("POST /api/refresh", cfg.RefreshToken)
	mux.HandleFunc("POST /api/revoke", cfg.RefreshRevoke)
	mux.HandleFunc("PUT /api/users", cfg.Authentication)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.DeleteChirp)
	mux.HandleFunc("POST /api/polka/webhooks", cfg.WebHook)
	err = srvr.ListenAndServe()
	if err != nil {
		log.Fatalf("couldn't listen and serve from server - %v", err)
	}
}
