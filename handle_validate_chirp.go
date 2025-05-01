package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func HandleValidateChirp(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Body string `json:"body"`
	}
	input := Input{}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Printf("error decoding params: %v", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if len(input.Body) > 140 {
		response := map[string]string{"error": "Chirp too long"}

		data, err := json.Marshal(response)
		if err != nil {
			log.Printf("Error encoding json error response: %v", err)
			w.WriteHeader(500)
			return
		}
		w.WriteHeader(400)
		w.Write(data)
		return
	}

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
	response := map[string]string{"cleaned_body": strings.Join(clean_words, " ")}

	data, err := json.Marshal(response)
	if err != nil {
		log.Printf("Error encoding json error response: %v", err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(200)
	w.Write(data)

}
