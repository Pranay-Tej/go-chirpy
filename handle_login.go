package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Pranay-Tej/go-chirpy/internal/auth"
)

func (apiConfig *ApiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	input := LoginSignupInput{}
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

	user := UserResponse{
		ID:        dbUser.ID,
		CreatedAt: dbUser.CreatedAt,
		UpdatedAt: dbUser.UpdatedAt,
		Email:     dbUser.Email,
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
