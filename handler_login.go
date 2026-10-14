package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/Ravinder2102/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds int    `json:"expires_in_seconds"`
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	decoder := json.NewDecoder(r.Body)
	request := requestBody{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode request", err)
		return
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), request.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	ok, err := auth.CheckPasswordHash(request.Password, user.HashedPassword)
	if err != nil || !ok {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	// Make jwt for user
	expirationTime := time.Hour
	if request.ExpiresInSeconds > 0 && request.ExpiresInSeconds < 3600 {
		expirationTime = time.Duration(request.ExpiresInSeconds) * time.Second
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(expirationTime))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create JWT for user", err)
		return
	}
	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: accessToken,
	})
}
