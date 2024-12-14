package main

import (
	"time"

	"github.com/Cprakhar/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID uuid.UUID `json:"id"`
	Email string `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time	`json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func databaseUserToUser(user database.User) User {
	return User{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func databaseChirpToChirp(chirp database.Chirp) Chirp {
	return Chirp{
		ID: chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body: chirp.Body,
		UserID: chirp.UserID,
	}
}

func databaseChirpsToChirps(databaseChirps []database.Chirp) []Chirp {
	chirps := make([]Chirp, 0)
	for _, chirp := range databaseChirps {
		chirps = append(chirps, databaseChirpToChirp(chirp))
	}
	return chirps
}