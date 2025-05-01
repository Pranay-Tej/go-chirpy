package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func (apiConfig *ApiConfig) handleGetChirpById(w http.ResponseWriter, r *http.Request) {
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
	chirpJson, err := json.Marshal(Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID.UUID,
	})
	if err != nil {
		log.Printf("error encoding json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(chirpJson)
}
