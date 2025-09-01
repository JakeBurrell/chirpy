package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/JakeBurrell/chirpy/internal/auth"
)

type TokenReponse struct {
	Token string `json:"token"`
}

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	// Validate token
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJson(w, http.StatusUnauthorized, errorResponse{
			Error: fmt.Sprintf("Failed to validate authorization: %v", err),
		})
		return
	}

	user_id, err := cfg.db.GetUserFromRefreshToken(r.Context(), token)
	if err != nil {
		respondWithJson(w, http.StatusUnauthorized, errorResponse{
			Error: "Invalid token",
		})
		return
	}

	accessToken, err := auth.MakeJWT(user_id, cfg.secret, time.Hour)
	if err != nil {
		log.Printf("failed to create jwt token: %v", err)
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Failed to create token",
		})
	}

	respondWithJson(w, http.StatusOK, TokenReponse{
		Token: accessToken,
	})
}
