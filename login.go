package main

import (
	"encoding/json"
	"net/http"
	"time"

	"chirpy/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Password         string `json:"password"`
		Email            string `json:"email"`
		ExpiresInSeconds *int   `json:"expires_in_seconds"`
	}

	params := parameters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Couldn't decode parameters", err)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	match, err := auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil || !match {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	const maxExpiry = time.Hour
	expiresIn := maxExpiry
	if params.ExpiresInSeconds != nil {
		expiresIn = time.Duration(*params.ExpiresInSeconds) * time.Second
		if expiresIn > maxExpiry {
			expiresIn = maxExpiry
		}
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtTokenSecret, expiresIn)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create token", err)
		return
	}

	type response struct {
		User
		Token string `json:"token"`
	}

	respondWithJSON(w, http.StatusOK, response{
		User: User{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
			Email:     user.Email,
		},
		Token: token,
	})
}
