package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	apiCfg := apiConfig{}

	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	fileserver := http.FileServer(http.Dir("."))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileserver)))
	serveMux.HandleFunc("GET /api/metrics", apiCfg.handleMetricsInc)
	serveMux.HandleFunc("POST /api/reset", apiCfg.handleMetricsReset)
	serveMux.HandleFunc("GET /api/healthz", handleReadiness)

	fmt.Println("Server running at port 8080")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server not starting")
	}
}
