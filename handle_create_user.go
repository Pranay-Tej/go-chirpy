package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (apiConfig *ApiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Email string `json:"email"`
	}
	input := Input{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("error decoding params %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := apiConfig.db.CrateUser(r.Context(), input.Email)
	if err != nil {
		log.Printf("error decoding params %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userJson := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	data, err := json.Marshal(userJson)
	if err != nil {
		log.Printf("error decoding params %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
