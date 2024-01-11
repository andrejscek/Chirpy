package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
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

	headers := strings.Split(r.Header.Get("Authorization"), " ")
	if len(headers) < 2 {
		RespondWithError(w, 401, "Unauthorized")
		return
	}
	token_string := headers[1]

	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(token_string, &claims, func(token *jwt.Token) (interface{}, error) { return []byte(cfg.jwtSecret), nil })
	if err != nil {
		RespondWithError(w, 401, "Unauthorized")
		return
	}
	if !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) || claims.Issuer != "chirpy-access" {
		RespondWithError(w, 401, "Unauthorized")
		return
	}

	if len(chirp.Body) > 140 {

		msg := fmt.Sprintf("Chirp too long: %s", chirp.Body)
		RespondWithError(w, 400, msg)
		return
	}

	author_id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	resp, err := cfg.db.CreateChirp(replaceProfanity(chirp.Body), author_id)
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	RespondWithJSON(w, 201, resp)

}

func (cfg *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {

	author_id := 0

	if len(r.URL.Query().Get("author_id")) > 0 {
		var err error
		author_id, err = strconv.Atoi(r.URL.Query().Get("author_id"))
		if err != nil {
			RespondWithError(w, 400, "Something went wrong")
			return
		}
	}

	chirps, err := cfg.db.GetChirps(author_id)
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

func (cfg *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	chirp_id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		RespondWithError(w, 404, "Invalid Chirp ID")
		return
	}

	type ResponseBody struct {
		Status string `json:"status"`
	}

	headers := strings.Split(r.Header.Get("Authorization"), " ")
	if len(headers) < 2 {
		RespondWithError(w, 401, "Unauthorized")
		return
	}

	token_string := headers[1]

	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(token_string, &claims, func(token *jwt.Token) (interface{}, error) { return []byte(cfg.jwtSecret), nil })
	if err != nil {
		RespondWithError(w, 401, "Unauthorized")
		return
	}

	if !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) || claims.Issuer != "chirpy-access" {
		RespondWithError(w, 401, "Unauthorized")
		return
	}

	author_id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	err = cfg.db.DeleteChirp(author_id, chirp_id)
	if err != nil {
		if err.Error() == "forbidden" {
			RespondWithError(w, 403, "Forbidden")
			return
		}

		if err.Error() == "chirp not found" {
			RespondWithError(w, 404, "Chirp not found")
			return
		}

		RespondWithError(w, 400, "Something went wrong")
		return
	}

	RespondWithJSON(w, 200, ResponseBody{Status: "ok"})
}
