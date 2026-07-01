package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	js, err := json.Marshal(payload)
	if err != nil {
		log.Fatal(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(js)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type Response struct {
		Error string `json:"error"`
	}
	respErr := Response{
		Error: msg,
	}

	respondWithJSON(w, code, respErr)
}
