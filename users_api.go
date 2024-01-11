package main

import (
	"Chirpy/database"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) createUser(w http.ResponseWriter, r *http.Request) {
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

		resp, err := cfg.db.CreateUser(params.Email, hash)
		if err != nil {
			RespondWithError(w, 400, "Something went wrong when creating user")
			return
		}

		RespondWithJSON(w, 201, UserResponse{ID: resp.ID, Email: resp.Email})
	} else {
		RespondWithError(w, 400, "Missing Email or Password")
	}
}

func (cfg *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Expires  int    `json:"expires_in_seconds,omitempty"`
	}

	type UserResponse struct {
		ID      int    `json:"id"`
		Email   string `json:"email"`
		Token   string `json:"token"`
		Refresh string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		RespondWithError(w, 400, "Something went wrong when decoding parameters")
		return
	}

	if (len(params.Email) > 0) && (len(params.Password) > 0) {
		user, err := cfg.db.GetUser(params.Email)
		if err != nil {
			RespondWithError(w, 401, "User not found")
			return
		}

		err = bcrypt.CompareHashAndPassword(user.Password, []byte(params.Password))
		if err != nil {
			RespondWithError(w, 401, "Wrong password")
			return
		}

		access_claims := jwt.RegisteredClaims{
			Issuer:    "chirpy-access",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(1) * time.Hour)),
			Subject:   fmt.Sprintf("%d", user.ID),
		}

		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, access_claims).SignedString([]byte(cfg.jwtSecret))
		if err != nil {
			RespondWithError(w, 400, "Something went wrong when creating token")
			return
		}

		refresh_claims := jwt.RegisteredClaims{
			Issuer:    "chirpy-refresh",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(24*60) * time.Hour)),
			Subject:   fmt.Sprintf("%d", user.ID),
		}

		refresh, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refresh_claims).SignedString([]byte(cfg.jwtSecret))
		if err != nil {
			RespondWithError(w, 400, "Something went wrong")
			return
		}

		RespondWithJSON(w, 200, UserResponse{ID: user.ID, Email: user.Email, Token: token, Refresh: refresh})
	} else {
		RespondWithError(w, 400, "Missing Email or Password")
	}
}

func (cfg *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	type Parameters struct {
		Pwd   string `json:"password"`
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

	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		RespondWithError(w, 500, "Something went wrong")
	}

	if (len(params.Email) > 0) && (len(params.Pwd) > 0) {
		hash, err := bcrypt.GenerateFromPassword([]byte(params.Pwd), 4)
		if err != nil {
			RespondWithError(w, 400, "Something went wrong")
		}

		resp, err := cfg.db.UpdateUser(id, params.Email, hash)
		if err != nil {
			RespondWithError(w, 400, "Something went wrong")
			return
		}

		RespondWithJSON(w, 200, UserResponse{ID: resp.ID, Email: resp.Email})
	} else {
		RespondWithError(w, 400, "Something went wrong")
	}
}

func (cfg *apiConfig) refreshAccess(w http.ResponseWriter, r *http.Request) {
	type TokenResponse struct {
		Token string `json:"token"`
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

	if !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) || claims.Issuer != "chirpy-refresh" {
		RespondWithError(w, 401, "Unauthorized")
		return
	}

	revoked, err := cfg.db.CheckRevoked(token_string)
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}
	if revoked != (database.Token{}) {
		RespondWithError(w, 401, "Unauthorized")
		return
	}

	access_claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(1) * time.Hour)),
		Subject:   claims.Subject,
	}

	access, err := jwt.NewWithClaims(jwt.SigningMethodHS256, access_claims).SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	RespondWithJSON(w, 200, TokenResponse{Token: access})
}

func (c *apiConfig) revokeRefresh(w http.ResponseWriter, r *http.Request) {
	type RevokedResponse struct {
		Revoked    string    `json:"revoked_token"`
		RevokeTime time.Time `json:"revoked_at"`
	}

	headers := strings.Split(r.Header.Get("Authorization"), " ")
	if len(headers) < 2 {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	token_string := headers[1]

	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(token_string, &claims, func(token *jwt.Token) (interface{}, error) { return []byte(c.jwtSecret), nil })
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	if !token.Valid || claims.ExpiresAt.Time.Before(time.Now()) || claims.Issuer != "chirpy-refresh" {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	revoked, err := c.db.RevokeToken(token_string)
	if err != nil {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	if revoked == (database.Token{}) {
		RespondWithError(w, 400, "Something went wrong")
		return
	}

	RespondWithJSON(w, 200, RevokedResponse{Revoked: revoked.ID, RevokeTime: revoked.RevokeTime})
}
