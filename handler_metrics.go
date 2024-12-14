package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) handleMetricsInc(w http.ResponseWriter, r *http.Request) {
	hits := apiCfg.fileserverHits.Load()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hits: %d", hits)
}

func (apiCfg *apiConfig) handleMetricsReset(w http.ResponseWriter, r *http.Request) {
	apiCfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits counter reset to 0"))
}