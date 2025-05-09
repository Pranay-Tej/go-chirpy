package main

import (
	"encoding/json"
	"log"
	"net/http"
	"slices"

	"github.com/Pranay-Tej/go-chirpy/internal/database"
	"github.com/google/uuid"
)

func (apiConfig *ApiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("author_id")
	authorId, err := uuid.Parse(id)
	if id != "" && err != nil {
		log.Printf("error parsing author_id: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var chirps []database.Chirp

	if authorId != uuid.Nil {
		chirps, err = apiConfig.db.GetChirpsByAuthorId(r.Context(), uuid.NullUUID{
			UUID:  authorId,
			Valid: true,
		})
		if err != nil {
			log.Printf("error fetching author chirps: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	} else {
		chirps, err = apiConfig.db.GetAllChirps(r.Context())
		if err != nil {
			log.Printf("error fetching chirps, %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
	sorting := r.URL.Query().Get("sort")
	if sorting == "desc" {
		slices.Reverse(chirps)
	}

	var jsonChirps []Chirp = make([]Chirp, 0)
	for _, chirp := range chirps {
		jsonChirps = append(jsonChirps, Chirp{
			Id:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID.UUID,
		})
	}

	data, err := json.Marshal(jsonChirps)
	if err != nil {
		log.Printf("error encoding json: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
