package main

import (
	"encoding/json"
	"net/http"
)

func (c *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Email string `json:"email"`
	}

	type UserResponse struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	if len(params.Email) > 0 {
		resp, err := c.db.CreateUser(params.Email)
		if err != nil {
			RespondWithError(w, 400, "Something went wrong")
			return
		}

		RespondWithJSON(w, 201, UserResponse{ID: resp.ID, Email: resp.Email})
	} else {
		RespondWithError(w, 400, "Something went wrong")
	}
}
