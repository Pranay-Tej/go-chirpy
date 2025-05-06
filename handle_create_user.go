package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Pranay-Tej/go-chirpy/internal/auth"
	"github.com/Pranay-Tej/go-chirpy/internal/database"
)

func (apiConfig *ApiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	input := SignupInput{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("error decoding params %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
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
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
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
