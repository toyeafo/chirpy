package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/toyeafo/chirpy/internal/database"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	Body      string    `json:"body"`
}

func (cfg *apiConfig) handleChirps(wr http.ResponseWriter, r *http.Request) {
	type req_body struct {
		Body   string    `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	req_body_text := req_body{}
	err := decoder.Decode(&req_body_text)
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError,
			"Couldn't decode request params", err)
		return
	}

	const maxChirpLength = 140
	if len(req_body_text.Body) > maxChirpLength {
		respondwithError(wr, http.StatusBadRequest,
			"Chirp length longer than 140 characters", nil)
		return
	}

	cleaned_body := cleanChirp(req_body_text.Body)

	chirp, err := cfg.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   cleaned_body,
		UserID: req_body_text.UserID,
	})
	if err != nil {
		log.Fatalf("error creating user in database: %s", err)
		wr.WriteHeader(500)
	}

	respondwithJSON(wr, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      cleaned_body,
		UserID:    chirp.UserID,
	})
}
