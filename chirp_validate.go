package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"strings"
)

func handlerChirpValidate(w http.ResponseWriter, r *http.Request) {
	const sizeLimit int = 140

	type requestParams struct {
		Body string `json:"body"`
	}

	type validResponse struct {
		CleanedBody string `json:"cleaned_body"`
	}

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
	if len(params.Body) > sizeLimit {
		respondWithJson(w, http.StatusBadRequest, errorResponse{
			Error: "Chrip is too long",
		})
		return
	}
	respondWithJson(w, http.StatusOK, validResponse{
		CleanedBody: replaceProfanity(params.Body),
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
