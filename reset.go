package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {

	if cfg.platform != "dev" {
		respondWithJson(w, http.StatusForbidden, errorResponse{
			Error: "Something went wrong",
		})
		log.Printf("Attempt to delete users outside of development")
		return
	}

	err := cfg.db.DeleteUsers(r.Context())
	if err != nil {
		respondWithJson(w, http.StatusForbidden, errorResponse{
			Error: "Something went wrong",
		})
		log.Printf("Failed to delete users from database: %v", err)
	}
	log.Printf("All users deleted from database")

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits rest to 0 and all users deleted"))
}
