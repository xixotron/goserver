package main

import (
	"errors"
	"net/http"
	"sort"

	"github.com/google/uuid"
)

func (cfg *apiConfig) handleGetChirps(w http.ResponseWriter, r *http.Request) {
	response := []Chirp{}

	chirps, err := cfg.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "unnable to get chirps", err)
		return
	}

	authorID := uuid.Nil
	authorIDstring := r.URL.Query().Get("author_id")
	if authorIDstring != "" {
		authorID, err = validateUUID(&authorIDstring)
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Malformed author id", err)
			return
		}
	}

	order := r.URL.Query().Get("sort")
	if order != "" && order != "asc" && order != "desc" {
		respondWithError(w, http.StatusBadRequest, "sort should be 'asc' or 'desc' if present", err)
		return
	}
	if order == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	for _, chirp := range chirps {
		if authorID != uuid.Nil && authorID != chirp.UserID {
			continue
		}
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
