package main

import (
	"strings"
)

func WordValidation(chirp string) string {
	invalidWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	if len(chirp) == 0 {
		return ""
	}
	var chirpymod []string
	s := strings.Fields(chirp)
	for _, word := range s {
		if _, ok := invalidWords[strings.ToLower(word)]; !ok {
			chirpymod = append(chirpymod, word)
		} else {
			chirpymod = append(chirpymod, "****")
		}
	}
	newchirpy := strings.Join(chirpymod, " ")
	return newchirpy
}
