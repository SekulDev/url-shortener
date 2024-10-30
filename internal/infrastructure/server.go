package infrastructure

import (
	"github.com/go-redis/redis"
	"log"
	"os"
	"strconv"
	"url-shortener/internal/app/service"
	"url-shortener/internal/app/usecase"
	"url-shortener/internal/domain/repository"
	mongoDb "url-shortener/internal/infrastructure/database/mongo"
	redisDb "url-shortener/internal/infrastructure/database/redis"
	"url-shortener/pkg"
)

type Repositories struct {
	UrlRepository *repository.MongoUrlRepository
}

type Usecases struct {
	HashUsecase *usecase.HashUsecase
	UrlUsecase  *usecase.UrlUsecaseImpl
}

type Services struct {
	TemplateService *service.TemplateService
	UrlService      *service.UrlServiceImpl
}

type Server struct {
	Mongo        *mongoDb.MongoClient
	Redis        *redisDb.RedisClient
	Repositories *Repositories
	Usecases     *Usecases
	Services     *Services
}

func NewServer() *Server {
	mongoClient, err := mongoDb.NewMongoClient(os.Getenv("MONGO_URL"), os.Getenv("MONGO_DATABASE"))
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	redisClient := redisDb.NewRedisClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	var snowflakeID int64
	snowflakeEnv := os.Getenv("SNOWFLAKE_NODE_ID")
	if snowflakeEnv == "" {
		snowflakeID = 1
	} else {
		converted, convertErr := strconv.ParseInt(snowflakeEnv, 10, 64)
		if convertErr != nil {
			log.Fatalf("Failed to convert snowflake ID: %v", convertErr)
		}
		snowflakeID = converted
	}

	// snowflake node
	node := pkg.InitSnowflakeNode(snowflakeID)

	// repositories
	urlRepository := repository.NewMongoUrlRepository(mongoClient.GetDatabase())

	// usecases
	hashUsecase := usecase.NewHashUsecase(node)
	urlUsecase := usecase.NewUrlUsecase(urlRepository, redisClient)
	rateLimitUsecase := usecase.NewRateLimitUsecase(redisClient)

	// services
	templateService := service.NewTemplateService()
	urlService := service.NewUrlService(urlUsecase, hashUsecase, rateLimitUsecase)

	server := &Server{
		Mongo: mongoClient,
		Redis: redisClient,
		Repositories: &Repositories{
			UrlRepository: urlRepository,
		},
		Usecases: &Usecases{
			HashUsecase: hashUsecase,
			UrlUsecase:  urlUsecase,
		},
		Services: &Services{
			TemplateService: templateService,
			UrlService:      urlService,
		},
	}

	return server
}
