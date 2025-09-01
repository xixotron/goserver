package main

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/xixotron/goserver/internal/database"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	response := []Chirp{}

	var err error
	authorID := uuid.Nil
	authorIDstring := r.URL.Query().Get("author_id")
	if authorIDstring != "" {
		authorID, err = validateUUID(&authorIDstring)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Malformed author id", err)
			return
		}
	}

	var chirps []database.Chirp
	if authorID != uuid.Nil {
		chirps, err = cfg.db.GetChirpsByUserID(r.Context(), authorID)
	} else {
		chirps, err = cfg.db.GetAllChirps(r.Context())
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unnable to get chirps", err)
		return
	}

	for _, chirp := range chirps {
		response = append(response, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (cfg *apiConfig) handleGetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDstring := r.PathValue("chirpID")
	chirpID, err := validateUUID(&chirpIDstring)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}

	chirp, err := cfg.db.GetChirpByID(r.Context(), chirpID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			respondWithError(w, http.StatusNotFound, "couldn't get chirp", err)
			return
		}
		respondWithError(w, http.StatusInternalServerError, "something went wrogn", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})
}

func validateUUID(id *string) (uuid.UUID, error) {
	if id == nil {
		return uuid.Nil, errors.New("no uuid provided")
	}
	resultID, err := uuid.Parse(*id)
	if err != nil {
		return uuid.Nil, errors.New("invalid user_id provided")
	}
	return resultID, nil
}
