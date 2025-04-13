package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/toyeafo/chirpy/internal/auth"
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
		Body string `json:"body"`
		// UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	req_body_text := req_body{}
	err := decoder.Decode(&req_body_text)
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError,
			"Couldn't decode request params", err)
		return
	}

	headerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "no valid tokens", err)
		return
	}

	userID, err := auth.ValidateJWT(headerToken, cfg.secret)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "invalid token", err)
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
		UserID: userID,
	})
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError, "Couldn't create chirp", err)
		return
	}

	respondwithJSON(wr, 201, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      cleaned_body,
		UserID:    chirp.UserID,
	})
}

func (cfg *apiConfig) handleChirpsGet(wr http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetChirps(r.Context())
	if err != nil {
		log.Fatalf("error getting chirps from database: %s", err)
		wr.WriteHeader(500)
	}

	var chirpList []Chirp
	for _, val := range chirps {
		chirpList = append(chirpList, Chirp{
			ID:        val.ID,
			CreatedAt: val.CreatedAt,
			UpdatedAt: val.UpdatedAt,
			Body:      val.Body,
			UserID:    val.UserID,
		})
	}
	respondwithJSON(wr, http.StatusOK, chirpList)

}

func (cfg *apiConfig) handleChirpGetSingle(wr http.ResponseWriter, r *http.Request) {
	idVal, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondwithError(wr, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetSingleChirp(r.Context(), idVal)
	if err != nil {
		respondwithError(wr, http.StatusBadRequest, "error retrieving chirp", err)
		return
	}

	respondwithJSON(wr, http.StatusOK, Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	})

}
