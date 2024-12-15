package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Cprakhar/Chirpy/internal/auth"
	"github.com/Cprakhar/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Webhook struct {
	Event string `json:"event"`
	Data struct {
		UserId uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (apiCfg *apiConfig) handleWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("No API key found: %v", err))
		return
	}

	if apiKey != apiCfg.polkaKey {
		responseWithError(w, 401, "Invalid API key")
		return
	}

	var webhook Webhook
	err = json.NewDecoder(r.Body).Decode(&webhook)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse the JSON: %v", err))
		return
	}
	if webhook.Event != "user.upgraded" {
		responseWithError(w, 204, "Error serving the request")
		return
	}
	_, err = apiCfg.dbQueries.UpgradeUserToRed(r.Context(), database.UpgradeUserToRedParams{
		UpdatedAt: time.Now().UTC(),
		ID: webhook.Data.UserId,
	})
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("User not in the database: %v", err))
		return
	}
	responseWithJSON(w, 204, struct{}{})
}