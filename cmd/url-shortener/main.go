package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"url-shortener/internal/adapters"
	"url-shortener/internal/infrastructure/database/mongo"
	"url-shortener/internal/infrastructure/hashing"
)

func main() {
	mongoClient, err := database.NewMongoClient(os.Getenv("MONGO_URL"), os.Getenv("MONGO_DATABASE"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect()

	//@TODO make it everything in config bootstrap
	node := hashing.InitSnowflakeNode(1)
	hashService := hashing.NewHashService(node)
	id := hashService.GenerateHash()
	fmt.Printf("Base62  ID: %s\n", id)

	//_mongoRepository := repository.NewMongoUrlRepository(mongoClient.GetDatabase())

	router := adapters.NewRouter()
	http.Handle("/", router)

	err = http.ListenAndServe(os.Getenv("HTTP_PORT"), nil)
	if err != nil {
		log.Fatal(err)
	}
}
