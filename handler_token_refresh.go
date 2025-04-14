package main

import (
	"net/http"
	"time"

	"github.com/toyeafo/chirpy/internal/auth"
)

func (cfg *apiConfig) handleTokenRefresh(wr http.ResponseWriter, req *http.Request) {
	headerToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "no valid tokens", err)
		return
	}

	user, err := cfg.db.GetUserRefreshToken(req.Context(), headerToken)
	if err != nil {
		respondwithError(wr, 401, "could not find token for user", err)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.secret, time.Hour)
	if err != nil {
		respondwithError(wr, http.StatusInternalServerError, "error creating token", err)
		return
	}

	respondwithJSON(wr, http.StatusOK, User{Token: token})
}

func (cfg *apiConfig) handleTokenRevoke(wr http.ResponseWriter, req *http.Request) {
	refreshToken, err := auth.GetBearerToken(req.Header)
	if err != nil {
		respondwithError(wr, http.StatusUnauthorized, "no valid tokens", err)
		return
	}

	err = cfg.db.UpdateRefreshToken(req.Context(), refreshToken)
	if err != nil {
		respondwithError(wr, 401, "error revoking token in db", err)
		return
	}
	wr.WriteHeader(http.StatusNoContent)
}
