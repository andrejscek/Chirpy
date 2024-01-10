package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	type ErrorBody struct {
		Error string `json:"error"`
	}

	log.Printf("Error: %s", msg)

	error := ErrorBody{
		Error: msg,
	}

	resp, err := json.Marshal(error)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}

func RespondChirpBody(w http.ResponseWriter, msg string) {
	type ChirpBody struct {
		Body string `json:"body"`
	}

	error := ChirpBody{
		Body: msg,
	}

	resp, err := json.Marshal(error)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(resp)
	log.Printf("Chirp validated: %s", msg)

}
