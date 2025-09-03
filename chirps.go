package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/JakeBurrell/chirpy/internal/auth"
	"github.com/JakeBurrell/chirpy/internal/database"
	"github.com/google/uuid"
)

type chirpJson struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Invalid chirp id: %v", err)
		respondWithJson(w, http.StatusNotFound, errorResponse{
			Error: "Invalid chirp id provided",
		})
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		log.Printf("Failed to retrieve chirp: %v", err)
		respondWithJson(w, http.StatusNotFound, errorResponse{
			Error: "Failed to retrieve chirp",
		})
		return
	}
	respondWithJson(w, http.StatusOK, chirpJson{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}

func (cfg *apiConfig) handlerAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("Could not retrieve chirps: %v", err)
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Failed to retrieve chirps from database",
		})
		return
	}

	jsonChirps := []chirpJson{}
	for _, chirp := range chirps {
		jsonChirps = append(jsonChirps, chirpJson{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	respondWithJson(w, http.StatusOK, jsonChirps)

}

func (cfg *apiConfig) handlerDeleteChirps(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Invalid chirp id: %v", err)
		respondWithJson(w, http.StatusBadRequest, errorResponse{
			Error: "Invalid chirp id provided",
		})
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJson(w, http.StatusUnauthorized, errorResponse{
			Error: fmt.Sprintf("Failed to validate authorization: %v", err),
		})
		return
	}

	loggedInID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithJson(w, http.StatusUnauthorized, errorResponse{
			Error: fmt.Sprintf("Invalid token: %v", err),
		})
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		respondWithJson(w, http.StatusNotFound, errorResponse{
			Error: "Chirp does not exist",
		})
		return
	}

	if chirp.UserID != loggedInID {
		respondWithJson(w, http.StatusForbidden, errorResponse{
			Error: "You are not the owner of this chirp",
		})
		return
	}

	err = cfg.db.DeleteChirp(r.Context(), chirpID)
	if err != nil {
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Something went wrong",
		})
		return
	}

	respondWithJson(w, http.StatusNoContent, nil)
	log.Printf("Chirp %s was deleted", chirpID)

}

func (cfg *apiConfig) handlerCreateChirps(w http.ResponseWriter, r *http.Request) {
	const sizeLimit int = 140

	type requestParams struct {
		Body string `json:"body"`
	}

	// Decode Body
	decoder := json.NewDecoder(r.Body)
	params := requestParams{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Something went wrong",
		})
		return
	}

	// Validate user
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithJson(w, http.StatusUnauthorized, errorResponse{
			Error: fmt.Sprintf("Failed to validate authorization: %v", err),
		})
		return
	}
	loggedInID, err := auth.ValidateJWT(token, cfg.secret)
	if err != nil {
		respondWithJson(w, http.StatusUnauthorized, errorResponse{
			Error: err.Error(),
		})
		return
	}

	if len(params.Body) > sizeLimit {
		respondWithJson(w, http.StatusBadRequest, errorResponse{
			Error: "Chrip is too long",
		})
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   replaceProfanity(params.Body),
		UserID: loggedInID,
	})
	if err != nil {
		log.Printf("Failed to add chirp to database: %v for user %s", err, loggedInID)
		respondWithJson(w, http.StatusInternalServerError, errorResponse{
			Error: "Couldn't add chrip to database",
		})
		return
	}
	log.Printf("New chirp Created")
	respondWithJson(w, http.StatusCreated, chirpJson{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    loggedInID,
	})

}

func replaceProfanity(chirp string) string {
	profanities := []string{"kerfuffle", "sharbert", "fornax"}
	cleaned := []string{}
	for _, word := range strings.Split(chirp, " ") {
		lowWord := strings.ToLower(word)
		if slices.Contains(profanities, lowWord) {
			cleaned = append(cleaned, "****")
		} else {
			cleaned = append(cleaned, word)
		}
	}
	return strings.Join(cleaned, " ")
}
