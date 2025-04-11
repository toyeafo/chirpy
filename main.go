package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/toyeafo/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
}

func main() {
	const pathRoot = "."
	const port = "8080"

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB URL not set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Printf("error opening a connection to the database: %s", err)
	}
	defer dbConn.Close()
	dbQueries := database.New(dbConn)

	mux := http.NewServeMux()
	server := &http.Server{Handler: mux, Addr: ":" + port}
	apiCfg := &apiConfig{fileserverHits: atomic.Int32{}, db: dbQueries}

	filehandler := apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", http.FileServer(http.Dir(pathRoot))))

	mux.Handle("/app/", filehandler)
	mux.HandleFunc("GET /api/healthz", handleHealth)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleHits)
	mux.HandleFunc("GET /api/chirps", apiCfg.handleChirpsGet)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleChirpGetSingle)
	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)
	mux.HandleFunc("POST /api/chirps", apiCfg.handleChirps)
	server.ListenAndServe()
}
