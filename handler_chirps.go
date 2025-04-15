package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
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

func (cfg *apiConfig) handleChirps(wr http.ResponseWriter, req *http.Request) {
	type req_body struct {
		Body string `json:"body"`
		// UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(req.Body)
	req_body_text := req_body{}
	err := decoder.Decode(&req_body_text)
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError,
			"Couldn't decode request params", err)
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

	const maxChirpLength = 140
	if len(req_body_text.Body) > maxChirpLength {
		respondwithError(wr, http.StatusBadRequest,
			"Chirp length longer than 140 characters", nil)
		return
	}

	cleaned_body := cleanChirp(req_body_text.Body)

	chirp, err := cfg.db.CreateChirp(req.Context(), database.CreateChirpParams{
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

func (cfg *apiConfig) handleChirpsGet(wr http.ResponseWriter, req *http.Request) {
	chirps, err := cfg.db.GetChirps(req.Context())
	if err != nil {
		log.Fatalf("error getting chirps from database: %s", err)
		wr.WriteHeader(500)
	}

	sortBy := req.URL.Query().Get("sort")

	authorID := uuid.Nil
	authorIDHeader := req.URL.Query().Get("author_id")
	if authorIDHeader != "" {
		authorID, err = uuid.Parse(authorIDHeader)
		if err != nil {
			respondwithError(wr, http.StatusNotFound, "error parsing authors id", err)
			return
		}
	}

	var chirpList []Chirp
	for _, val := range chirps {
		if authorID != uuid.Nil && val.UserID != authorID {
			continue
		}
		chirpList = append(chirpList, Chirp{
			ID:        val.ID,
			CreatedAt: val.CreatedAt,
			UpdatedAt: val.UpdatedAt,
			Body:      val.Body,
			UserID:    val.UserID,
		})
	}

	if sortBy == "desc" {
		sort.Slice(chirpList, func(i, j int) bool {
			return chirpList[i].CreatedAt.After(chirpList[j].CreatedAt)
		})
	}
	respondwithJSON(wr, http.StatusOK, chirpList)

}

func (cfg *apiConfig) handleChirpGetSingle(wr http.ResponseWriter, req *http.Request) {
	idVal, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondwithError(wr, http.StatusBadRequest, "Invalid chirp ID", err)
		return
	}
	chirp, err := cfg.db.GetSingleChirp(req.Context(), idVal)
	if err != nil {
		respondwithError(wr, http.StatusNotFound, "error retrieving chirp", err)
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

func (cfg *apiConfig) handleChirpDelete(wr http.ResponseWriter, req *http.Request) {
	chirpID, err := uuid.Parse(req.PathValue("chirpID"))
	if err != nil {
		respondwithError(wr, 401, "Invalid chirp ID", err)
		return
	}

	accessToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "missing or invalid access token", err)
		return
	}

	user, err := auth.ValidateJWT(accessToken, cfg.secret)
	if err != nil {
		respondwithError(wr, 403, "invalid access token", err)
		return
	}

	_, err = cfg.db.GetSingleChirpByIDandUser(
		req.Context(),
		database.GetSingleChirpByIDandUserParams{
			UserID: user,
			ID:     chirpID,
		},
	)
	if err != nil {
		respondwithError(wr, http.StatusForbidden, "chirp can't be found", err)
		return
	}

	err = cfg.db.DeleteChirpByID(req.Context(), database.DeleteChirpByIDParams{
		UserID: user,
		ID:     chirpID,
	})
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "user not allowed to delete chirp", err)
		return
	}
	wr.WriteHeader(204)
}
