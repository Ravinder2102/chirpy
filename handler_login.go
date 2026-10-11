package main

import (
	"encoding/json"
	"net/http"

	"github.com/Ravinder2102/chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}

	type response struct {
		User
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
	respondWithJson(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
	})
}
