package main

import "net/http"

func (cfg *apiConfig) handleReset(wr http.ResponseWriter, req *http.Request) {
	cfg.fileserverHits.Store(0)
	wr.Write([]byte("Hits resets to 0."))
}
