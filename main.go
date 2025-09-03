package main

import (
	"database/sql"
	"github.com/JakeBurrell/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secret         string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	// Connect to databse
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platformEnv := os.Getenv("PLATFORM")
	secretEnv := os.Getenv("SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	dbQueries := database.New(db)

	cfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platformEnv,
		secret:         secretEnv,
	}

	mux := http.NewServeMux()
	mux.Handle(
		"/app/",
		http.StripPrefix(
			"/app",
			cfg.middlewareMetricsInc(http.FileServer(http.Dir(filepathRoot))),
		),
	)
	mux.HandleFunc("GET /api/healthz", handlerHealthz)
	mux.HandleFunc("GET /admin/metrics", cfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", cfg.handlerReset)
	mux.HandleFunc("POST /api/users", cfg.handlerCreateUser)
	mux.HandleFunc("POST /api/chirps", cfg.handlerCreateChirps)
	mux.HandleFunc("POST /api/login", cfg.handleLogin)
	mux.HandleFunc("POST /api/refresh", cfg.handlerRefresh)
	mux.HandleFunc("POST /api/revoke", cfg.handlerRevoke)
	mux.HandleFunc("GET /api/chirps", cfg.handlerAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.handlerGetChirp)
	mux.HandleFunc("PUT /api/users", cfg.handlerUpdateUser)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.handlerDeleteChirps)

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(server.ListenAndServe())
}
