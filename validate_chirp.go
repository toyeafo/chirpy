package main

import (
	"encoding/json"
	"net/http"
)

func handleValidateChirp(wr http.ResponseWriter, r *http.Request) {
	type req_body struct {
		Body string `json:"body"`
	}

	type resp_body struct {
		Body string `json:"cleaned_body"`
	}

	decoder := json.NewDecoder(r.Body)
	req_body_text := req_body{}
	err := decoder.Decode(&req_body_text)
	if err != nil {
		respondwithError(wr, http.StatusBadRequest, "Couldn't decode request params", err)
		return
	}

	const maxChirpLength = 140
	if len(req_body_text.Body) > maxChirpLength {
		respondwithError(wr, http.StatusBadRequest, "Chirp length longer than 140 characters", nil)
		return
	}

	cleaned_body := cleanChirp(req_body_text.Body)

	respondwithJSON(wr, http.StatusOK, resp_body{
		Body: cleaned_body,
	})
}
