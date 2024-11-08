package main

import (
	"log"
	"net/http"
	"os"
	"url-shortener/internal/adapters"
	"url-shortener/internal/infrastructure"
)

func main() {
	server := infrastructure.NewServer()

	router := adapters.NewRouter(server)
	http.Handle("/", router)

	log.Printf("Starting server on port %s", os.Getenv("HTTP_PORT"))
	err := http.ListenAndServe(":"+os.Getenv("HTTP_PORT"), nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
		return
	}
}
