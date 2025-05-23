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
	polkaKey := os.Getenv("POLKA_KEY")
	if polkaKey == "" {
		log.Fatal("POLKA_KEY env not set")
	}
	apiConfig := ApiConfig{
		db:        dbQueries,
		platform:  platform,
		jwtSecret: jwtSecret,
		polkaKey:  polkaKey,
	}

	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("POST /admin/reset", apiConfig.handleReset)
	mux.HandleFunc("GET /admin/metrics", apiConfig.handleMetrics)
	mux.HandleFunc("POST /api/users", apiConfig.handleCreateUser)
	mux.HandleFunc("PUT /api/users", apiConfig.handleUpdateUser)
	mux.HandleFunc("POST /api/login", apiConfig.handleLogin)
	mux.HandleFunc("POST /api/refresh", apiConfig.handleRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiConfig.handleRevokeRefreshToken)
	mux.HandleFunc("GET /api/healthz", HandleHealthz)
	mux.HandleFunc("POST /api/chirps", apiConfig.handleCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.handleGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{id}", apiConfig.handleGetChirpById)
	mux.HandleFunc("DELETE /api/chirps/{id}", apiConfig.handleDeleteChirpById)
	mux.HandleFunc("POST /api/polka/webhooks", apiConfig.handlePolkaWebhook)

	log.Printf("Serving on port: %s\n", port)

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
