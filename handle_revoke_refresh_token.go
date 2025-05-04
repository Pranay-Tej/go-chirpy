package main

import (
	"log"
	"net/http"

	"github.com/Pranay-Tej/go-chirpy/internal/auth"
)

func (apiConfig *ApiConfig) handleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("bad bearer token: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = apiConfig.db.RevokeRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf("unable to revoke token: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
