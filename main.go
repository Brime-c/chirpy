package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	var apiCfg apiConfig
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))
	prefixed := http.StripPrefix("/app", fileServer)
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(prefixed))
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /api/validate_chirp", func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Body string `json:"body"`
		}
		decoder := json.NewDecoder(r.Body)
		resp := response{}
		err := decoder.Decode(&resp)
		if err != nil {
			respondWithError(w, 500, "Something went wrong")
			return
		}
		if len(resp.Body) > 140 {
			respondWithError(w, 400, "Chirp is too long")
			return
		}
		cleanedBody := getCleanedBody(resp.Body)

		type returnVals struct {
			CleanedBody string `json:"cleaned_body"`
		}
		val := returnVals{
			CleanedBody: cleanedBody,
		}
		respondWithJSON(w, 200, val)
	})
	mux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)
	myServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(myServer.ListenAndServe())
}
