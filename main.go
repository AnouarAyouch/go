package main

import (
	"fmt"
	"net/http"
)

func start() {
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	start()
}

