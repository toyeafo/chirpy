package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/toyeafo/chirpy/internal/auth"
	"github.com/toyeafo/chirpy/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
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
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})

}
