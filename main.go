package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/arturogood17/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	hits      atomic.Int32
	dbQueries *database.Queries
}

func main() {
	godotenv.Load(".env")
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)
	const port = "8080"
	mux := http.NewServeMux()
	h := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	a := apiConfig{}
	a.dbQueries = dbQueries
	mux.Handle("/app/", a.middlewareMetricsInc(h)) //Hay que quitar el app porque queremos servir
	//los archivos que están en el dir actual
	mux.HandleFunc("GET /api/healthz", handlerReadiness) //no tienes que crear un directorio para el path
	mux.HandleFunc("GET /admin/metrics", a.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", a.HandlerReset)
	mux.HandleFunc("POST /api/validate_chirp", HandlerChirps)
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
