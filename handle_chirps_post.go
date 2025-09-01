package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/xixotron/goserver/internal/auth"
	"github.com/xixotron/goserver/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlePostChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body *string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer token not provided or invalid", err)
		return
	}

	userID, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer token not provided or invalid", err)
		return
	}
	user, err := cfg.db.GetUserByID(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer token not provided or invalid", err)
		return
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err = decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusTeapot, "Couldn't parse provided data", err)
		return
	}

	chirpBody, err := validateChirp(params.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error(), err)
		return
	}

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		ID:     uuid.New(),
		Body:   chirpBody,
		UserID: user.ID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError,
			"Couldn't create chirp", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validateChirp(chirp *string) (string, error) {
	const maxChirpLength = 140

	if chirp == nil {
		return "", errors.New("chirp body not provided")
	}

	if *chirp == "" {
		return "", errors.New("empty chirp body provided")
	}

	if len(*chirp) > maxChirpLength {
		return "", errors.New("Chirp is too long")
	}

	return replaceNotyWords(*chirp), nil
}

func replaceNotyWords(chirp string) string {
	notyWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	res := chirp
	for word := range strings.SplitSeq(res, " ") {
		_, match := notyWords[strings.ToLower(word)]
		if match {
			res = strings.ReplaceAll(res, word, "****")
		}
	}
	return res
}
