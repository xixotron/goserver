package main

import (
	"encoding/json"
	"net/http"
)

func handleChirpsValidate(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140

	type parameters struct {
		Body *string `json:"body"`
	}

	type returnVals struct {
		Valid bool `json:"valid"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusTeapot, "Couldn't parse provided data", err)
		return
	}

	if params.Body == nil {
		respondWithError(w, http.StatusBadRequest, "No body parameter provided", nil)
		return
	}

	if len(*params.Body) > maxChirpLength {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long", nil)
		return
	}

	respondWithJSON(w, http.StatusOK, returnVals{Valid: true})
}
