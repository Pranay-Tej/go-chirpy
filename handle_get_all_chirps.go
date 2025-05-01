package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func (apiConfig *ApiConfig) handleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := apiConfig.db.GetAllChirps(r.Context())
	if err != nil {
		log.Printf("error fetching chirps, %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
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
