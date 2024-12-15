package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Cprakhar/Chirpy/internal/auth"
	"github.com/Cprakhar/Chirpy/internal/database"
	"github.com/google/uuid"
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
	defaultExpiresInSeconds := time.Hour.Seconds()
	tokenString, err := auth.MakeJWT(logUser.ID, apiCfg.tokenSecret, time.Duration(defaultExpiresInSeconds)*time.Second)
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error logging user: %v", err))
		return
	}
	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		responseWithError(w, 401, fmt.Sprintf("Error logging user: %v", err))
		return
	}
	_, err = apiCfg.dbQueries.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token: refreshToken,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID: logUser.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 60),
		RevokedAt: sql.NullTime{},
	})
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't save the refresh token in the database: %v", err))
		return
	}
	
	loggedUser := LoggedUser{
		ID: logUser.ID,
		CreatedAt: logUser.CreatedAt,
		UpdatedAt: logUser.UpdatedAt,
		Email: logUser.Email,
		Token: tokenString,
		RefreshToken: refreshToken,
	}
	responseWithJSON(w, 200, loggedUser)
}

func (apiCfg *apiConfig) handleUpdateLoginDetails (w http.ResponseWriter, r *http.Request) {
	var user LoginUser
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		responseWithError(w, 400, fmt.Sprintf("Couldn't parse the JSON: %v", err))
		return
	}
	defer r.Body.Close()

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
	hashedPasswd, err := auth.HashPassword(user.Password)
	if err != nil {
		responseWithError(w, 500, fmt.Sprintf("Couldn't change the password: %v", err))
		return
	}

	updatedUser, err := apiCfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		Email: user.Email,
		HashedPassword: hashedPasswd,
		UpdatedAt: time.Now().UTC(),
		ID: userId,
	})
	if err != nil {
		responseWithError(w, 404, fmt.Sprintf("Couldn't found the user: %v", err))
		return
	}
	responseWithJSON(w, 200, databaseUserToUser(updatedUser))
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
