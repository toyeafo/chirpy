package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/toyeafo/chirpy/internal/auth"
	"github.com/toyeafo/chirpy/internal/database"
)

type User struct {
	ID               uuid.UUID `json:"id"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Email            string    `json:"email"`
	Token            string    `json:"token"`
	RefreshToken     string    `json:"refresh_token"`
	ChirpyMembership bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handleCreateUser(wr http.ResponseWriter, req *http.Request) {
	type req_body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	req_body_text := req_body{}
	err := json.NewDecoder(req.Body).Decode(&req_body_text)
	if err != nil {
		respondwithError(wr, http.StatusBadRequest, "Couldn't decode request params", err)
		return
	}

	secret_pwd, err := auth.HashPassword(req_body_text.Password)
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError, "error securing password", err)
		return
	}

	user, err := cfg.db.CreateUser(req.Context(), database.CreateUserParams{
		Email:          req_body_text.Email,
		HashedPassword: secret_pwd,
	})
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError, "error creating user in database", err)
		return
	}

	respondwithJSON(wr, 201, User{
		ID:               user.ID,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		Email:            user.Email,
		ChirpyMembership: user.IsChirpyRed,
	})

}

func (cfg *apiConfig) handleUserUpdate(wr http.ResponseWriter, req *http.Request) {
	type req_body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	req_body_text := req_body{}
	err := json.NewDecoder(req.Body).Decode(&req_body_text)
	if err != nil {
		respondwithError(wr, http.StatusBadRequest, "Couldn't decode request params", err)
		return
	}

	secret_pwd, err := auth.HashPassword(req_body_text.Password)
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError, "error securing password", err)
		return
	}

	headerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "no valid tokens", err)
		return
	}

	userID, err := auth.ValidateJWT(headerToken, cfg.secret)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "invalid token", err)
		return
	}

	user, err := cfg.db.UpdateUsernamePassword(req.Context(), database.UpdateUsernamePasswordParams{
		Email:          req_body_text.Email,
		HashedPassword: secret_pwd,
		ID:             userID,
	})
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError, "error updating user in database", err)
		return
	}

	respondwithJSON(wr, http.StatusOK, User{
		ID:               user.ID,
		CreatedAt:        user.CreatedAt,
		UpdatedAt:        user.UpdatedAt,
		Email:            user.Email,
		Token:            headerToken,
		ChirpyMembership: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlePolkaWebhook(wr http.ResponseWriter, req *http.Request) {
	type UserData struct {
		UserID uuid.UUID `json:"user_id"`
	}

	type req_body struct {
		Event string   `json:"event"`
		Data  UserData `json:"data"`
	}

	apiPolka, err := auth.GetAPIKey(req.Header)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "api key not found", err)
		return
	}

	if cfg.polka_key != apiPolka {
		respondwithError(wr, http.StatusUnauthorized, "invalid api key", nil)
		return
	}

	req_body_text := req_body{}
	err = json.NewDecoder(req.Body).Decode(&req_body_text)
	if err != nil {
		respondwithError(wr, http.StatusBadRequest, "Couldn't decode request params", err)
		return
	}

	if req_body_text.Event != "user.upgraded" {
		respondwithJSON(wr, 204, "")
		return
	}

	_, err = cfg.db.UpgradeUser(req.Context(), database.UpgradeUserParams{
		IsChirpyRed: true,
		ID:          req_body_text.Data.UserID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondwithError(wr, http.StatusNotFound, "Couldn't find user", err)
			return
		}
		respondwithError(wr, http.StatusNotFound, "could not upgrade user", err)
		return
	}

	respondwithJSON(wr, 204, "")
}
