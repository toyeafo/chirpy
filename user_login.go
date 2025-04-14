package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/toyeafo/chirpy/internal/auth"
	"github.com/toyeafo/chirpy/internal/database"
)

func (cfg *apiConfig) handleUserLogin(wr http.ResponseWriter, req *http.Request) {
	type req_body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		// Expiry   time.Duration `json:"expires_in_seconds"`
	}

	req_body_text := req_body{}
	err := json.NewDecoder(req.Body).Decode(&req_body_text)
	if err != nil {
		respondwithError(wr, http.StatusBadRequest, "Couldn't decode request params", err)
		return
	}

	user, err := cfg.db.RetrieveUserPwd(req.Context(), req_body_text.Email)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(req_body_text.Password, user.HashedPassword)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError, "error creating token", err)
		return
	}

	refreshToken, err := auth.MakeRefreshToken()
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError, "error creating refresh token", err)
		return
	}

	_, err = cfg.db.CreateRefreshToken(req.Context(),
		database.CreateRefreshTokenParams{
			Token:     refreshToken,
			UserID:    user.ID,
			ExpiresAt: time.Now().AddDate(0, 0, 60),
		})
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError, "error saving refresh token to db", err)
		return
	}

	respondwithJSON(wr, 200, User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
