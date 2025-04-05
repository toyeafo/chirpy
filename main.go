package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	pathRoot := "."
	const port = "8080"

	mux := http.NewServeMux()
	server := &http.Server{Handler: mux, Addr: ":" + port}

	filehandler := http.FileServer(http.Dir(pathRoot))
	apiCfg := apiConfig{}

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app/", filehandler)))
	mux.HandleFunc("/healthz", handleHealth)
	mux.HandleFunc("/metrics", apiCfg.handleHits)
	mux.HandleFunc("/reset", apiCfg.handleReset)
	server.ListenAndServe()
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleHits(wr http.ResponseWriter, req *http.Request) {
	hitresponse := fmt.Sprintf("Hits: %d", cfg.fileserverHits.Load())
	wr.Write([]byte(hitresponse))
}
