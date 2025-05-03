package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/Pranay-Tej/go-chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	const port = "8080"

	mux := http.NewServeMux()
	httpServer := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	godotenv.Load()

	platform := os.Getenv("PLATFORM")
	if platform == "" {
		log.Fatal("PLATFORM env not set")
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL env not set")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("error connection to db")
	}
	dbQueries := database.New(db)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET env not set")
	}
	apiConfig := ApiConfig{
		db:        dbQueries,
		platform:  platform,
		jwtSecret: jwtSecret,
	}

	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("POST /admin/reset", apiConfig.handleReset)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handleMetrics)
	mux.HandleFunc("POST /api/users", apiConfig.handleCreateUser)
	mux.HandleFunc("POST /api/login", apiConfig.handleLogin)
	mux.HandleFunc("GET /api/healthz", HandleHealthz)
	mux.HandleFunc("POST /api/chirps", apiConfig.handleCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.handleGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiConfig.handleGetChirpById)

	log.Printf("Serving on port: %s\n", port)

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
