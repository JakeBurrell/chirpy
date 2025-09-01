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

type LoginRequest struct {
	Password string `json:"password"`
	Email    string `json:"email"`
}

type LoginResponse struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
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

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("failed to create jwt token: %v", err)
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Failed to create token",
		})
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		log.Printf("failed to create refresh token")
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Failed to create token",
		})
	}

	err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:  refreshToken,
		UserID: user.ID,
	})
	if err != nil {
		log.Printf("failed to add refresh token to database")
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Failed to create token",
		})
	}

	respondWithJson(w, http.StatusOK, LoginResponse{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshToken,
	})

}
