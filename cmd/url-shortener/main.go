package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"os"
	"url-shortener/internal/adapters"
	mongoDb "url-shortener/internal/infrastructure/database/mongo"
	redisDb "url-shortener/internal/infrastructure/database/redis"
	"url-shortener/internal/infrastructure/hashing"
)

func main() {
	mongoClient, err := mongoDb.NewMongoClient(os.Getenv("MONGO_URL"), os.Getenv("MONGO_DATABASE"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect()

	redisClient := redisDb.NewRedisClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	defer redisClient.Disconnect()

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
