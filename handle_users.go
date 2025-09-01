package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/xixotron/goserver/internal/auth"
	"github.com/xixotron/goserver/internal/database"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password *string `json:"password"`
		Email    *string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	if params.Email == nil || *params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "email parameter was not provided or empty", nil)
		return
	}
	if params.Password == nil || *params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "password parameter was not provided or empty", nil)
		return
	}
	passworHash, err := auth.HashPassword(*params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "password parameter invalid", nil)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          *params.Email,
		HashedPassword: passworHash,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handleModifyUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password *string `json:"password"`
		Email    *string `json:"email"`
	}

	accessToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Bearer token not provided or invalid", err)
		return
	}

	userID, err := auth.ValidateJWT(accessToken, cfg.jwtSecret)
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
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	if params.Email == nil || *params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "email parameter was not provided or empty", nil)
		return
	}

	if params.Password == nil || *params.Password == "" {
		respondWithError(w, http.StatusBadRequest, "password parameter was not provided or empty", nil)
		return
	}

	passworHash, err := auth.HashPassword(*params.Password)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "password parameter invalid", nil)
		return
	}

	user, err = cfg.db.UpdateUserCredentials(r.Context(), database.UpdateUserCredentialsParams{
		ID:             user.ID,
		Email:          *params.Email,
		HashedPassword: passworHash,
	})

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, User{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}
