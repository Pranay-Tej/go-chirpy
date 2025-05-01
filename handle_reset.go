package main

import (
	"fmt"
	"net/http"
)

func (apiConfig *ApiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if apiConfig.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	apiConfig.fileServerHits.Store(0)
	err := apiConfig.db.DeleteAllUsers(r.Context())
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
