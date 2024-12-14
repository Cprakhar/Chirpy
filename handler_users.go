package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Cprakhar/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/Cprakhar/Chirpy/internal/auth"
)

type LoginUser struct {
	Email string `json:"email"`
	Password string `json:"password"`
}

func (apiCfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse the JSON: %v", err))
		return
	}
	defer r.Body.Close()

	passwd, err := auth.HashPassword(user.Password)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Password doesn't meet the requirements: %v", err))
		return
	}
	
	newUser, err := apiCfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		Email: user.Email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		HashedPassword: passwd,
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't create a user: %v", err))
		return
	}

	responseWithJSON(w, 201, databaseUserToUser(newUser))
}

func (apiCfg *apiConfig) handleUserLogin(w http.ResponseWriter, r *http.Request) {
	var user LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse the JSON: %v", err))
		return
	}
	logUser, err := apiCfg.dbQueries.GetUserByEmail(r.Context(), user.Email)
	if err != nil {
		responseWithError(w, 401, "Incorrect email or password")
		return
	}
	err = auth.CheckHashPassword(user.Password, logUser.HashedPassword)
	if err != nil {
		responseWithError(w, 401, "Incorrect email or password")
		return
	}
	responseWithJSON(w, 200, databaseUserToUser(logUser))
}


func (apiCfg *apiConfig) handleDeleteAllUsers(w http.ResponseWriter, r *http.Request){
	if apiCfg.platform != "dev" {
		responseWithError(w, 403, "Forbidden in production to delete all users")
		return
	}
	err := apiCfg.dbQueries.DeleteAllUsers(r.Context())
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't delete all users: %v", err))
	}
	responseWithJSON(w, 200, struct{}{})
}
