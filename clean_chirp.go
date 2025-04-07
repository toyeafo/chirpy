package main

import "strings"

func cleanChirp(chirp string) string {
	profanewords := [3]string{"kerfuffle", "sharbert", "fornax"}
	split_body_text := strings.Split(chirp, " ")
	for ind, word := range split_body_text {
		for _, profaneword := range profanewords {
			val := strings.Compare(strings.ToLower(word), strings.ToLower(profaneword))
			if val == 0 {
				split_body_text[ind] = strings.Repeat("*", 4)
			}
		}
	}

	return strings.Join(split_body_text, " ")
}
