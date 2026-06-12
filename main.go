package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	myServer := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	err := myServer.ListenAndServe()
	if err != nil {
		fmt.Errorf("error while starting the server")
	}
}
