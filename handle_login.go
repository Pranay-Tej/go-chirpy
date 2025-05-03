package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Pranay-Tej/go-chirpy/internal/auth"
	"github.com/google/uuid"
)

type LoginResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
}

type LoginInput struct {
	Email            string `json:"email"`
	Password         string `json:"password"`
	ExpiresInSeconds int64  `json:"expires_in_seconds"`
}

func (apiConfig *ApiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	input := LoginInput{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("error decoding params %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dbUser, err := apiConfig.db.GetUserByEmail(r.Context(), input.Email)
	if err != nil {
		log.Printf("user not found: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	err = auth.CheckPasswordHash(dbUser.HashedPassword, input.Password)
	if err != nil {
		log.Printf("Incorrect email or password: %v", err)
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		mismatchError := map[string]string{"error": "Incorrect email or password"}
		res, err := json.Marshal(mismatchError)
		if err != nil {
			log.Printf("error encoding json%v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(res)
		return
	}

	fmt.Printf("input.ExpiresInSeconds %v", input.ExpiresInSeconds)
	expiresIn := time.Hour
	if input.ExpiresInSeconds > 0 && input.ExpiresInSeconds < int64(time.Hour.Seconds()) {
		expiresIn = time.Duration(input.ExpiresInSeconds) * time.Second
	}
	token, err := auth.MakeJwt(dbUser.ID, apiConfig.jwtSecret, expiresIn)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user := LoginResponse{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
		Token:     token,
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		log.Printf("error encoding json%v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(userJson)
}
