package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func respondWithJson(w http.ResponseWriter, code int, jsonResponse any) {
	w.Header().Set("Content-Type", "application/json")
	dat, err := json.Marshal(jsonResponse)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		return
	}
	w.WriteHeader(code)
	w.Write(dat)
}
