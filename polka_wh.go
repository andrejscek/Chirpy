package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (cfg *apiConfig) poklaWebhook(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Event string `json:"event"`
		Data  struct {
			UserID int `json:"user_id"`
		} `json:"data"`
	}

	type ResponseBody struct {
		Status string `json:"status"`
	}

	headers := strings.Split(r.Header.Get("Authorization"), " ")
	if len(headers) < 2 || headers[1] != cfg.polkaKey {
		RespondWithError(w, 401, "Unauthorized")
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	if params.Event == "user.upgraded" {
		err = cfg.db.UpgradeUser(params.Data.UserID)
		if err != nil {
			RespondWithError(w, 404, "user not found")
			return
		}
	}

	RespondWithJSON(w, 200, ResponseBody{})
}
