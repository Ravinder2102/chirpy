package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/Ravinder2102/chirpy/internal/database"
	"github.com/google/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	type response struct {
		Post
	}

	decoder := json.NewDecoder(r.Body)
	request := requestBody{}
	err := decoder.Decode(&request)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode request", err)
		return
	}

	const maxChirpLength = 140
	if len(request.Body) > maxChirpLength {
		respondWithError(w, 400, "Chirp is too long", nil)
		return
	}

	cleanedBody := replaceProfaneWords(request.Body)

	post, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleanedBody,
		UserID: request.UserID})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't create post", err)
		return
	}

	respondWithJson(w, http.StatusCreated, response{
		Post: Post{
			ID:        post.ID,
			CreatedAt: post.CreatedAt,
			UpdatedAt: post.UpdatedAt,
			Body:      post.Body,
			UserID:    post.UserID,
		},
	})

}

func replaceProfaneWords(text string) string {
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	words := strings.Split(text, " ")
	for i, word := range words {
		if _, exist := profaneWords[strings.ToLower(word)]; exist {
			words[i] = "****"
		}
	}

	cleanedText := strings.Join(words, " ")
	return cleanedText
}
