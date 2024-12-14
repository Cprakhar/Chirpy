package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Cprakhar/Chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpBody struct {
	Body string `json:"body"`
	UserId uuid.UUID `json:"user_id"`
}

func (apiCfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {

	var chirp chirpBody

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&chirp)
	if err != nil {
		responseWithError(w, 400, "Something went wrong")
		return
	}

	defer r.Body.Close()

	if chirp.Body == "" {
		responseWithError(w, 400, "Something went wrong")
		return
	}

	if len(chirp.Body) > 140 {
		responseWithError(w, 400, "Chirp is too long")
		return
	}

	cleanedChirp := replaceProfane(chirp.Body)

	newChirp, err := apiCfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Body: cleanedChirp,
		UserID: chirp.UserId,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't create a chirp: %v", err))
		return
	}

	responseWithJSON(w, 201, databaseChirpToChirp(newChirp))
}

func (apiCfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := apiCfg.dbQueries.GetAllChirps(r.Context())
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Couldn't retrieve all chirps: %v", err))
	}
	responseWithJSON(w, 200, databaseChirpsToChirps(chirps))
}

func (apiCfg *apiConfig) handleGetChirpByID(w http.ResponseWriter, r *http.Request) {
	chirpId := r.PathValue("chirpID")
	id, err := uuid.Parse(chirpId)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse the chirpId: %v", err))
		return
	}

	chirp, err := apiCfg.dbQueries.GetChirpByID(r.Context(), id)
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Couldn't get the chirp: %v", err))
		return
	}
	responseWithJSON(w, 200, databaseChirpToChirp(chirp))
}


func replaceProfane(chirp string) string {
	profane := map[string]bool{
		"kerfuffle": true,
		"sharbert": true,
		"fornax" : true,
	}

	chirpWords := strings.Split(chirp, " ")
	for i, chirpWord := range chirpWords {
		chirpWord = strings.ToLower(chirpWord)
		if _, ok := profane[chirpWord]; ok {
			chirpWords[i] = "****"
		}
	}
	return strings.Join(chirpWords, " ")
}