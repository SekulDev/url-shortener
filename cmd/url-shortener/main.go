package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"os"
	"url-shortener/internal/adapters"
	"url-shortener/internal/app/usecase"
	mongoDb "url-shortener/internal/infrastructure/database/mongo"
	redisDb "url-shortener/internal/infrastructure/database/redis"
	"url-shortener/pkg"
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
	node := pkg.InitSnowflakeNode(1)
	hashUsecase := usecase.NewHashUsecase(node)
	id := hashUsecase.GenerateHash()
	fmt.Printf("Base62  ID: %s\n", id)

	//mongoRepository := repository.NewMongoUrlRepository(mongoClient.GetDatabase())
	//uc := usecase.NewUrlUsecase(mongoRepository, redisClient)
	//us := service.NewUrlService(uc)
	//
	//resultUrl, err := us.ResolveShortUrl("2CeYV8b0XYn")
	//if err != nil {
	//	fmt.Printf("Failed to resolve short url: %v", err)
	//	return
	//}
	//fmt.Println(resultUrl.ShortId, resultUrl.LongUrl)

	router := adapters.NewRouter()
	http.Handle("/", router)

	http.ListenAndServe(os.Getenv("HTTP_PORT"), nil)
}
