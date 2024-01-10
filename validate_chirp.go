package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func replaceProfanity(msg string) string {
	words := strings.Split(msg, " ")

	bad_words := []string{"kerfuffle", "sharbert", "fornax"}

	for i, word := range words {
		for _, bad_word := range bad_words {
			if strings.ToLower(word) == bad_word {
				words[i] = "****"
			}
		}
	}

	return strings.Join(words, " ")
}

func (cfg *apiConfig) handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type Chirp struct {
		Body string `json:"body"`
	}
	decoder := json.NewDecoder(r.Body)
	chirp := Chirp{}
	err := decoder.Decode(&chirp)
	if err != nil {
		msg := fmt.Sprintf("Error decoding parameters: %s", err)
		RespondWithError(w, 500, msg)
		return
	}
	if len(chirp.Body) > 140 {
		msg := fmt.Sprintf("Chirp too long: %s", chirp.Body)
		RespondWithError(w, 400, msg)
		return
	}
	RespondWithJSON(w, 200, Chirp{Body: replaceProfanity(chirp.Body)})

}
