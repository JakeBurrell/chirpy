package main

import (
	"encoding/json"
	"github.com/google/uuid"
	"log"
	"net/http"
	"time"
)

func (cfg *apiConfig) handlersCreateUser(w http.ResponseWriter, r *http.Request) {

	type createUserJson struct {
		Email string `json:"email"`
	}

	type userJson struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

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

	user, err := cfg.db.CreateUser(r.Context(), params.Email)
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
