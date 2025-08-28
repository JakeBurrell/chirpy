package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/JakeBurrell/chirpy/internal/auth"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Password   string `json:"password"`
	Email      string `json:"email"`
	ExpiresSec *int   `json:"expires_in_seconds"`
}

type LoginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := LoginRequest{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding user parameters: %s", err)
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Couldn't decode parameters",
		})
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithJson(w, http.StatusBadRequest, errorResponse{
			Error: "User does not exist",
		})
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithJson(w, http.StatusUnauthorized, errorResponse{
			Error: "Incorrect email or password",
		})
		return
	}

	var expirationTime time.Duration
	if params.ExpiresSec != nil {
		expirationTime = min(
			time.Duration(*params.ExpiresSec)*time.Second,
			time.Hour,
		)
	} else {
		expirationTime = time.Hour
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, expirationTime)
	if err != nil {
		log.Printf("failed to create token: %v", err)
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Failed to create token",
		})
	}

	respondWithJson(w, http.StatusOK, LoginResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
		Token:     token,
	})

}
