package main

import (
	"encoding/json"
	"net/http"

	"github.com/Pranay-Tej/go-chirpy/internal/auth"
	"github.com/Pranay-Tej/go-chirpy/internal/database"
	"github.com/google/uuid"
)

type webhookInput struct {
	Event string `json:"event"`
	Data  struct {
		UserId uuid.UUID `json:"user_id"`
	} `json:"data"`
}

const USER_UPGRADED_EVENT = "user.upgraded"

func (apiConfig *ApiConfig) handlePolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetPolkaApiKey(r.Header)
	if err != nil || apiKey != apiConfig.polkaKey {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	input := webhookInput{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if input.Event != USER_UPGRADED_EVENT {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err = apiConfig.db.SetUserChirpRedStatus(r.Context(), database.SetUserChirpRedStatusParams{
		IsChirpyRed: true,
		ID:          input.Data.UserId,
	})

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)

}
