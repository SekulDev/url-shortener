package infrastructure

import (
	"github.com/go-redis/redis"
	"github.com/xinguang/go-recaptcha"
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

type Recaptcha struct {
	Public string
	Secret string
}

type Repositories struct {
	UrlRepository *repository.MongoUrlRepository
}

type Usecases struct {
	HashUsecase      *usecase.HashUsecase
	UrlUsecase       *usecase.UrlUsecaseImpl
	RateLimitUsecase *usecase.RateLimitUsecaseImpl
	RecaptchaUsecase *usecase.RecaptchaUsecaseImpl
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
	Recaptcha    *Recaptcha
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

	// recaptcha
	recaptchaKeys := &Recaptcha{
		Public: os.Getenv("RECAPTCHA_PUBLIC"),
		Secret: os.Getenv("RECAPTCHA_SECRET"),
	}

	recaptchaObj, recaptchaErr := recaptcha.NewWithSecert(recaptchaKeys.Secret)
	if recaptchaErr != nil {
		log.Fatalf("Failed to initialize recaptcha: %v", recaptchaErr)
	}

	// snowflake node
	node := pkg.InitSnowflakeNode(snowflakeID)

	// repositories
	urlRepository := repository.NewMongoUrlRepository(mongoClient.GetDatabase())

	// usecases
	hashUsecase := usecase.NewHashUsecase(node)
	urlUsecase := usecase.NewUrlUsecase(urlRepository, redisClient)
	rateLimitUsecase := usecase.NewRateLimitUsecase(redisClient)
	recaptchaUsecase := usecase.NewRecaptchaUsecase(recaptchaObj)

	// services
	templateService := service.NewTemplateService(recaptchaKeys.Public)
	urlService := service.NewUrlService(urlUsecase, hashUsecase, rateLimitUsecase, recaptchaUsecase)

	server := &Server{
		Mongo: mongoClient,
		Redis: redisClient,
		Repositories: &Repositories{
			UrlRepository: urlRepository,
		},
		Usecases: &Usecases{
			HashUsecase:      hashUsecase,
			UrlUsecase:       urlUsecase,
			RateLimitUsecase: rateLimitUsecase,
			RecaptchaUsecase: recaptchaUsecase,
		},
		Services: &Services{
			TemplateService: templateService,
			UrlService:      urlService,
		},
		Recaptcha: recaptchaKeys,
	}

	return server
}
