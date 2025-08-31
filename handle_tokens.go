package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/xixotron/goserver/internal/auth"
)

func (cfg *apiConfig) handleRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token not provided or invalid", fmt.Errorf("%v 1, %w", r.Pattern, err))
		return
	}

	err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "couldn't revoke the provided refresh token", fmt.Errorf("%v 2, %w", r.Pattern, err))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (cfg *apiConfig) handleRefreshToken(w http.ResponseWriter, r *http.Request) {
	type token struct {
		Token string `json:"token"`
	}

	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Refresh token not provided or invalid", fmt.Errorf("%v 1, %w", r.Pattern, err))
		return
	}

	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid or expired refresh token", fmt.Errorf("%v 2, %w", r.Pattern, err))
		return
	}

	bearerToken, err := auth.MakeJWT(userID, cfg.jwtSecret, tokenExpiration)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Coudn't create bearer token", fmt.Errorf("%v 3, %w", r.Pattern, err))
		return
	}
	log.Printf("%v 4: send bearer token\n", r.Pattern)
	respondWithJSON(w, http.StatusOK, token{
		Token: bearerToken,
	})
}
