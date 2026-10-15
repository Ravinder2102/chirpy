package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Ravinder2102/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateChirpyRed(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "couldn't find API Key", err)
		return
	}
	if apiKey != cfg.polkaKey {
		respondWithError(w, http.StatusUnauthorized, "API key is invalid", err)
		return
	}

	type requestBody struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	request := requestBody{}
	err = decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request", err)
		return
	}

	if request.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.db.UpdateChirpyRedByID(r.Context(), request.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "couldn't find user", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "couldn't update chirpy red status for user", err)
		return

	}
	w.WriteHeader(http.StatusNoContent)
}
