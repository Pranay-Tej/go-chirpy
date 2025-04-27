package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	fileServerHits atomic.Int32
}

func (apiConfig *ApiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiConfig.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (apiConfig *ApiConfig) reset(w http.ResponseWriter, r *http.Request) {
	apiConfig.fileServerHits.Store(0)
	w.WriteHeader(http.StatusOK)
}

func (apiConfig *ApiConfig) metrics(w http.ResponseWriter, r *http.Request) {
	hits := apiConfig.fileServerHits.Load()
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %v times!</p></body></html>", hits)))
}

func main() {
	const port = "8080"

	mux := http.NewServeMux()
	httpServer := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	apiConfig := ApiConfig{}
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.Handle("POST /admin/reset", http.HandlerFunc(apiConfig.reset))
	mux.Handle("GET /admin/metrics", http.HandlerFunc(apiConfig.metrics))

	mux.HandleFunc("GET /api/healthz", handleHealthz)

	log.Printf("Serving on port: %s\n", port)

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
