package main

import (
	"encoding/json"
	"net/http"

	"github.com/Ravinder2102/chirpy/internal/auth"
	"github.com/Ravinder2102/chirpy/internal/database"
)

func (cfg *apiConfig) handlerUpdateDetails(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	type requestBody struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type response struct {
		User
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token is missing or malformed", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "malformed or missing token", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	request := requestBody{}
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't decode request", err)
		return
	}

	hashedPassword, err := auth.HashPassword(request.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't hash password", err)
		return
	}

	user, err := cfg.db.UpdateUserEmailAndPassword(r.Context(), database.UpdateUserEmailAndPasswordParams{
		Email:          request.Email,
		HashedPassword: hashedPassword,
		ID:             userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't update details", err)
		return
	}

	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
