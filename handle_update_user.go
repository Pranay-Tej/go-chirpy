package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Pranay-Tej/go-chirpy/internal/auth"
	"github.com/Pranay-Tej/go-chirpy/internal/database"
)

func (apiConfig *ApiConfig) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	input := SignupInput{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("error decoding input: %v\n", err)
		w.WriteHeader(http.StatusBadRequest)
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

	hashedPassword, err := auth.HashPassword(input.Password)
	if err != nil {
		log.Printf("error hashing password%v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := apiConfig.db.UpdateUserCreds(r.Context(), database.UpdateUserCredsParams{
		Email:          input.Email,
		HashedPassword: hashedPassword,
		ID:             userId,
	})

	userJson := UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}

	data, err := json.Marshal(userJson)
	if err != nil {
		log.Printf("error encoding json: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)

}
