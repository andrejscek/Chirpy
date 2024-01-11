package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (c *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type UserResponse struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, "Something went wrong when decoding parameters")
		return
	}

	if (len(params.Email) > 0) && (len(params.Password) > 0) {
		hash, err := bcrypt.GenerateFromPassword([]byte(params.Password), 4)
		if err != nil {
			RespondWithError(w, 400, "Something went when hashing password")
		}

		resp, err := c.db.CreateUser(params.Email, hash)
		if err != nil {
			RespondWithError(w, 400, "Something went wrong when creating user")
			return
		}

		RespondWithJSON(w, 201, UserResponse{ID: resp.ID, Email: resp.Email})
	} else {
		RespondWithError(w, 400, "Missing Email or Password")
	}
}

func (c *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type UserResponse struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, "Something went wrong when decoding parameters")
		return
	}

	if (len(params.Email) > 0) && (len(params.Password) > 0) {
		user, err := c.db.GetUser(params.Email)
		if err != nil {
			RespondWithError(w, 401, "User not found")
			return
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
		if err != nil {
			RespondWithError(w, 401, "Wrong password")
			return
		}

		RespondWithJSON(w, 200, UserResponse{ID: user.ID, Email: user.Email})
	} else {
		RespondWithError(w, 400, "Missing Email or Password")
	}
}
