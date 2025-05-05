package main

import (
	"log"
	"net/http"

	"github.com/Pranay-Tej/go-chirpy/internal/auth"
	"github.com/google/uuid"
)

func (apiConfig *ApiConfig) handleDeleteChirpById(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	chirpId, err := uuid.Parse(id)
	if err != nil {
		log.Printf("error parsing id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	chirp, err := apiConfig.db.GetChirpById(r.Context(), chirpId)
	if err != nil {
		log.Printf("chirp not found: %v", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("no token found: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userId, err := auth.ValidateJwt(token, apiConfig.jwtSecret)
	if err != nil {
		log.Printf("token invalid: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if chirp.UserID.UUID != userId {
		log.Println("unauthrozied")
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = apiConfig.db.DeleteChirpById(r.Context(), chirp.ID)
	if err != nil {
		log.Printf("erro deleting chirp: %v\n", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
