package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
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

func (cfg *apiConfig) postChirp(w http.ResponseWriter, r *http.Request) {
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

	resp, err := cfg.db.CreateChirp(replaceProfanity(chirp.Body))
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	RespondWithJSON(w, 201, resp)

}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.db.GetChirps()
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	RespondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) getChirp(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondWithError(w, 400, "Invalid ID")
	}

	chirp, err := cfg.db.GetChirp(id)
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	RespondWithJSON(w, 200, chirp)
}
