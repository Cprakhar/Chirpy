package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/Cprakhar/Chirpy/internal/auth"
	"github.com/Cprakhar/Chirpy/internal/database"
)

func (apiCfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("No Authorization header was found: %v", err))
		return
	}
	
	refreshTokenDB, err := apiCfg.dbQueries.GetRefreshToken(r.Context(), refreshToken)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Refresh token not found: %v", err))
		return
	}
	if refreshTokenDB.ExpiresAt.Before(time.Now().UTC()) {
		responseWithError(w, 401, "Refresh token is expired")
		return
	}
	expiresIn := time.Hour.Seconds()
	tokenString, err := auth.MakeJWT(refreshTokenDB.UserID, apiCfg.tokenSecret, time.Duration(expiresIn)*time.Second)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Couldn't authorize the user: %v", err))
		return
	}

	responseWithJSON(w, 200, struct{Token string `json:"token"`}{Token: tokenString})
}

func (apiCfg *apiConfig) handleRevokeToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("No Authorization header was found: %v", err))
		return
	}

	err = apiCfg.dbQueries.UpdateRefreshToken(r.Context(), database.UpdateRefreshTokenParams{
		Token: refreshToken,
		UpdatedAt: time.Now().UTC(),
		RevokedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
	})
	
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Couldn't update the refresh token: %v", err))
		return
	}
	w.WriteHeader(204)
}