package main

import (
	"net/http"
	"os"
	"url-shortener/internal/adapters"
	"url-shortener/internal/infrastructure"
)

func main() {
	server := infrastructure.NewServer()

	router := adapters.NewRouter(server)
	http.Handle("/", router)

	http.ListenAndServe(os.Getenv("HTTP_PORT"), nil)
}
