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
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
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
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

const tokenExpiration = 1 * time.Hour

func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password *string `json:"password"`
		Email    *string `json:"email"`
	}

	type response struct {
		User
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
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

	user, err := cfg.db.GetUserByEmail(r.Context(), *params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(*params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	bearerToken, err := auth.MakeJWT(user.ID, cfg.jwtSecret, tokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coudn't create bearer token", err)
		return
	}

	refreshToken := auth.MakeRefreshToken()

	_, err = cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Hour * 24 * 60),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coudn't create refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token:        bearerToken,
		RefreshToken: refreshToken,
	})
}
