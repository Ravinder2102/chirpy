package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func handlerValidateChirps(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	type requestBody struct {
		Body string `json:"body"`
	}

	type responseBody struct {
		CleanedBody string `json:"cleaned_body"`
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

	respondWithJson(w, http.StatusOK, responseBody{CleanedBody: cleanedBody})

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
