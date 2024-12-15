package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/Cprakhar/Chirpy/internal/auth"
	"github.com/Cprakhar/Chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpBody struct {
	Body string `json:"body"`
}

func (apiCfg *apiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Unauthorized: %v", err))
		return
	}
	userId, err := auth.ValidateJWT(tokenString, apiCfg.tokenSecret)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("No authorized: %v", err))
		return
	}

	var chirp chirpBody

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&chirp)
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
		UserID: userId,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't create a chirp: %v", err))
		return
	}

	responseWithJSON(w, 201, databaseChirpToChirp(newChirp))
}

func (apiCfg *apiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	authorId := r.URL.Query().Get("author_id")
	sort := r.URL.Query().Get("sort")
	if authorId == "" {
		chirps, err := apiCfg.dbQueries.GetAllChirps(r.Context())
		if err != nil {
			responseWithError(w, 404, fmt.Sprintf("Couldn't retrieve all chirps: %v", err))
			return
		}
		sortedChirps(chirps, sort)
		responseWithJSON(w, 200, databaseChirpsToChirps(chirps))
		return
	}

	parsedAuthorId, err := uuid.Parse(authorId)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Error parsing: %v", err))
		return
	}
	chirps, err := apiCfg.dbQueries.GetAllChirpsByUserID(r.Context(), parsedAuthorId)
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("No user found: %v", err))
		return
	}
	sortedChirps(chirps, sort)
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

func (apiCfg *apiConfig) handleDeleteChirp(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Unauthorized: %v", err))
		return
	}
	userId, err := auth.ValidateJWT(tokenString, apiCfg.tokenSecret)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Unauthorized: %v", err))
		return
	}

	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse the uuid: %v", err))
		return
	}
	err = apiCfg.dbQueries.DeleteChirpByID(r.Context(), database.DeleteChirpByIDParams{
		UserID: userId,
		ID: chirpId,
	})
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("User not the owner or chirp not found: %v", err))
		return
	}
	responseWithJSON(w, 204, struct{}{})
}


func sortedChirps(chirps []database.Chirp, sortby string) {
	sort.Slice(chirps, func(i, j int) bool {
		if sortby == "asc" || sortby == "" {
			return chirps[i].CreatedAt.Before(chirps[j].CreatedAt)
		}
		return chirps[j].CreatedAt.Before(chirps[i].CreatedAt)
	})
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