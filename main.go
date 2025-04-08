package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/arturogood17/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	hits      atomic.Int32
	dbQueries *database.Queries
	PLATFORM  string
	SECRET    string
}

type User struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

type Chirp struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

func main() {
	const port = "8080"
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("db_url must be set")
	}
	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("Platform must be set")
	}
	secret := os.Getenv("SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET environment variable is not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening the database: %v", err)
	}
	dbQueries := database.New(db)
	a := apiConfig{
		hits:      atomic.Int32{},
		dbQueries: dbQueries,
		PLATFORM:  platform,
		SECRET:    secret,
	}
	mux := http.NewServeMux()
	h := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", a.middlewareMetricsInc(h)) //Hay que quitar el app porque queremos servir
	//los archivos que están en el dir actual
	mux.HandleFunc("GET /api/healthz", handlerReadiness) //no tienes que crear un directorio para el path
	mux.HandleFunc("GET /admin/metrics", a.HandlerMetrics)
	mux.HandleFunc("GET /api/chirps", a.AllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", a.SingleChirp)
	mux.HandleFunc("POST /admin/reset", a.HandlerReset)
	mux.HandleFunc("POST /api/users", a.HandlerUser)
	mux.HandleFunc("POST /api/chirps", a.HandlerChirps)
	mux.HandleFunc("POST /api/login", a.UserLogin)
	mux.HandleFunc("POST /api/refresh", a.RefreshToken)
	mux.HandleFunc("POST /api/revoke", a.Revoke)
	mux.HandleFunc("PUT /api/users", a.UpdateUser)
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	fmt.Printf("Serving files from on port: %s\n", port)
	log.Fatal(server.ListenAndServe())
}

func handlerReadiness(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.WriteHeader(http.StatusOK)
	res.Write([]byte("OK"))
}
