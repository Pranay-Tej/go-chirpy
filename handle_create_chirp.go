package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Pranay-Tej/go-chirpy/internal/auth"
	"github.com/Pranay-Tej/go-chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	Id        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserId    uuid.UUID `json:"user_id"`
}

func (apiConfig *ApiConfig) handleCreateChirp(w http.ResponseWriter, r *http.Request) {

	type Input struct {
		Body string `json:"body"`
	}
	input := Input{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("error decoding params: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	// TODO: make this util fn
	if len(input.Body) > 140 {
		response := map[string]string{"error": "Chirp too long"}

		data, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error encoding json error response: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Write(data)
		return
	}

	// TODO: make this util fn
	words := strings.Split(input.Body, " ")
	clean_words := make([]string, len(words))
	for i, word := range words {
		loweredWord := strings.ToLower(word)
		if loweredWord == "kerfuffle" || loweredWord == "sharbert" || loweredWord == "fornax" {
			clean_words[i] = "****"
			continue
		}
		clean_words[i] = word
	}
	cleanedBody := strings.Join(clean_words, " ")

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId, err := auth.ValidateJwt(token, apiConfig.jwtSecret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	chirp, err := apiConfig.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body: cleanedBody,
		UserID: uuid.NullUUID{
			UUID:  userId,
			Valid: true,
		},
	})

	if err != nil {
		log.Printf("Error creating chirp: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return

	}

	chirpJson := Chirp{
		Id:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID.UUID,
	}

	data, err := json.Marshal(chirpJson)
	if err != nil {
		log.Printf("Error encoding json error response: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Write(data)
}
