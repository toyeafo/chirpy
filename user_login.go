package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/toyeafo/chirpy/internal/auth"
)

func (cfg *apiConfig) handleUserLogin(wr http.ResponseWriter, req *http.Request) {
	type req_body struct {
		Email    string        `json:"email"`
		Password string        `json:"password"`
		Expiry   time.Duration `json:"expires_in_seconds"`
	}

	req_body_text := req_body{}
	err := json.NewDecoder(req.Body).Decode(&req_body_text)
	if err != nil {
		respondwithError(wr, http.StatusBadRequest, "Couldn't decode request params", err)
		return
	}

	var expiryTime time.Duration

	if req_body_text.Expiry == expiryTime {
		expiryTime = time.Hour
	} else if req_body_text.Expiry > time.Hour {
		expiryTime = time.Hour
	} else {
		expiryTime = req_body_text.Expiry
	}

	user_pwd, err := cfg.db.RetrieveUserPwd(req.Context(), req_body_text.Email)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(req_body_text.Password, user_pwd.HashedPassword)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	token, err := auth.MakeJWT(user_pwd.ID, cfg.secret, expiryTime)
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError, "error creating token", err)
		return
	}

	respondwithJSON(wr, 200, User{
		ID:        user_pwd.ID,
		CreatedAt: user_pwd.CreatedAt,
		UpdatedAt: user_pwd.UpdatedAt,
		Email:     user_pwd.Email,
		Token:     token,
	})
}
