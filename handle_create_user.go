package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Pranay-Tej/go-chirpy/internal/auth"
	"github.com/Pranay-Tej/go-chirpy/internal/database"
	"github.com/google/uuid"
)

type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}
type SignupInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (apiConfig *ApiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	input := SignupInput{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("error decoding params %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		log.Printf("error hashing password: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	user, err := apiConfig.db.CrateUser(r.Context(), database.CrateUserParams{
		Email:          input.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		log.Printf("error decoding params %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	userJson := UserResponse{
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
