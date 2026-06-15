package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))
	prefixed := http.StripPrefix("/app", fileServer)
	mux.Handle("/app/", prefixed)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	})

	myServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(myServer.ListenAndServe())

}
