package main

import "strings"

func getCleanedBody(s string) string {
	prof := []string{"kerfuffle", "sharbert", "fornax"}

	sString := strings.Split(s, " ")

	for i, word := range sString {
		for _, p := range prof {
			if strings.ToLower(word) == p {
				sString[i] = "****"
			}
		}
	}
	cleaned := strings.Join(sString, " ")
	return cleaned
}
