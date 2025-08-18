package main

import (
	"log"
	"net/http"
)

func FirstServer() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.Handle("./", http.FileServer(http.Dir("/assets")))

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	FirstServer()
}
