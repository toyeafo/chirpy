package main

import (
	"encoding/json"
	"net/http"

	"github.com/toyeafo/chirpy/internal/auth"
)

func (cfg *apiConfig) handleUserLogin(wr http.ResponseWriter, req *http.Request) {
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

	user_pwd, err := cfg.db.RetrieveUserPwd(req.Context(), req_body_text.Email)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	err = auth.CheckPasswordHash(user_pwd.HashedPassword, req_body_text.Password)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "Incorrect email or password", err)
		return
	}

	respondwithJSON(wr, 200, User{
		ID:        user_pwd.ID,
		CreatedAt: user_pwd.CreatedAt,
		UpdatedAt: user_pwd.UpdatedAt,
		Email:     user_pwd.Email,
	})
}
