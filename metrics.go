package main

import (
	"fmt"
	"net/http"
)

func (cfg *apiConfig) handlerNumberResponse(w http.ResponseWriter, r *http.Request) {
	count := cfg.fileserverHits.Load()
	response := fmt.Sprintf("Hits: %v", count)

	w.Write([]byte(response))
}

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Write([]byte("count reset"))
}
