package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/Cprakhar/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	platform string
	tokenSecret string
	polkaKey string
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Couldn't load the environment variables: %v\n", err)
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	tokenSecret := os.Getenv("TOKEN_SECRET")
	polkaKey := os.Getenv("POLKA_KEY")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Couldn't connect to the database: %v", err)
	}
	dbQueries := database.New(db)
	

	apiCfg := apiConfig{
		dbQueries: dbQueries,
		platform: platform,
		tokenSecret: tokenSecret,
		polkaKey: polkaKey,
	}

	serveMux := http.NewServeMux()
	server := &http.Server{
		Handler: serveMux,
		Addr:    ":8080",
	}

	fileserver := http.FileServer(http.Dir("."))
	serveMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileserver)))
	serveMux.HandleFunc("GET /admin/metrics", apiCfg.handleMetricsInc)
	serveMux.HandleFunc("POST /admin/reset", apiCfg.handleMetricsReset)
	serveMux.HandleFunc("GET /api/healthz", handleReadiness)

	serveMux.HandleFunc("GET /api/chirps", apiCfg.handleGetAllChirps)
	serveMux.HandleFunc("POST /api/chirps", apiCfg.handleCreateChirp)
	serveMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handleGetChirpByID)
	serveMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handleDeleteChirp)

	serveMux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)
	serveMux.HandleFunc("POST /api/login", apiCfg.handleUserLogin)
	serveMux.HandleFunc("PUT /api/users", apiCfg.handleUpdateLoginDetails)

	serveMux.HandleFunc("DELETE /admin/reset", apiCfg.handleDeleteAllUsers)
	serveMux.HandleFunc("POST /api/refresh", apiCfg.handleRefreshToken)
	serveMux.HandleFunc("POST /api/revoke", apiCfg.handleRevokeToken)

	serveMux.HandleFunc("POST /api/polka/webhooks", apiCfg.handleWebhook)


	fmt.Printf("Server running at port 8080, Platform: %s\n", platform)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Server not starting")
	}
}
