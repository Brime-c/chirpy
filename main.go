package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(".")))

	myServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal(myServer.ListenAndServe())

}
