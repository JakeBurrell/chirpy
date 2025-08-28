package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/JakeBurrell/chirpy/internal/auth"
	"github.com/JakeBurrell/chirpy/internal/database"
	"github.com/google/uuid"
)

type createUserJson struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type userJson struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlersCreateUser(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	params := createUserJson{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding user parameters: %s", err)
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Couldn't decode parameters",
		})
		return
	}

	password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			"Password could not be hashed",
		})
		return
	}
	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: password,
	})
	if err != nil {
		log.Printf("Error creating user in database: %s", err)
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Couldn't create user",
		})
	}

	respondWithJson(w, http.StatusCreated, userJson{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
	log.Printf("User Created: %v", user)

}
