package main

import (
	"log"
	"net/http"
	"os"
)

func (cfg *apiConfig) handleReset(wr http.ResponseWriter, req *http.Request) {
	platform := os.Getenv("PLATFORM")
	if platform != "dev" {
		log.Fatal("action not allowed")
		wr.WriteHeader(403)
	}
	cfg.fileserverHits.Store(0)
	cfg.db.DeleteUsers(req.Context())
	wr.Write([]byte("Hits resets to 0."))
}
