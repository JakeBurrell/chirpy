package main

import (
	"fmt"
	"net/http"

	"github.com/JakeBurrell/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJson(w, http.StatusBadRequest, errorResponse{
			Error: fmt.Sprintf("No token provided: %v", err),
		})
		return
	}

	err = cfg.db.RevokeToken(r.Context(), token)
	if err != nil {
		respondWithJson(w, http.StatusBadRequest, errorResponse{
			Error: "Invalid token",
		})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
