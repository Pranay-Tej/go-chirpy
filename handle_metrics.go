package main

import (
	"fmt"
	"net/http"
)

func (apiConfig *ApiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	hits := apiConfig.fileServerHits.Load()
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %v times!</p></body></html>", hits)))
}
