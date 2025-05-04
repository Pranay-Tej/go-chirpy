package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Pranay-Tej/go-chirpy/internal/auth"
)

func (apiConfig *ApiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("no bearer token: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	dbToken, err := apiConfig.db.GetRefreshToken(r.Context(), token)
	if err != nil {
		log.Printf("refresh token not found: %v\n", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if dbToken.RevokedAt.Valid {
		log.Println("refresh token is revoked")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if time.Until(dbToken.ExpiresAt) < 0 {
		log.Println("refresh token is expired")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newToken, err := auth.MakeJwt(dbToken.UserID, apiConfig.jwtSecret, time.Hour)
	if err != nil {
		log.Printf("error generating token: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	type Res struct {
		Token string `json:"token"`
	}
	res := Res{
		Token: newToken,
	}
	resJson, err := json.Marshal(res)
	if err != nil {
		log.Printf("error encoding json: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resJson)
}
