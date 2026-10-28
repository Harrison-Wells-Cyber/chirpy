package main

import (
	"net/http"
	"time"

	"chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "No token found", err)
		return
	}

	refreshToken, err := cfg.dbQueries.GetRefreshToken(r.Context(), tokenString)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Token not found", err)
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Token has been revoked", nil)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now().UTC()) {
		respondWithError(w, http.StatusUnauthorized, "Token has expired", nil)
		return
	}

	accessToken, err := auth.MakeJWT(refreshToken.UserID, cfg.jwtTokenSecret, time.Hour)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access token", err)
		return
	}

	type response struct {
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{Token: accessToken})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	tokenString, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't get token", err)
		return
	}

	if _, err := cfg.dbQueries.RevokeRefreshToken(r.Context(), tokenString); err != nil {
		respondWithError(w, http.StatusNotFound, "Token not found", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
