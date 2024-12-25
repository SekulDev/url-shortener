package infrastructure

import (
	"github.com/go-redis/redis"
	"github.com/xinguang/go-recaptcha"
	"log"
	"os"
	"strconv"
	"url-shortener/internal/app/service"
	urlS "url-shortener/internal/app/service/url"
	hashU "url-shortener/internal/app/usecase/hash"
	ratelimitU "url-shortener/internal/app/usecase/ratelimit"
	recaptchaU "url-shortener/internal/app/usecase/recaptcha"
	urlU "url-shortener/internal/app/usecase/url"
	repository "url-shortener/internal/domain/repository/url"
	mongoDb "url-shortener/internal/infrastructure/database/mongo"
	redisDb "url-shortener/internal/infrastructure/database/redis"
	mongoRepository "url-shortener/internal/infrastructure/repository/url/mongo"
	"url-shortener/pkg"
)

type Recaptcha struct {
	Public string
	Secret string
}

type Repositories struct {
	UrlRepository repository.UrlRepository
}

type Usecases struct {
	HashUsecase      hashU.HashUsecase
	UrlUsecase       urlU.UrlUsecase
	RateLimitUsecase ratelimitU.RatelimitUsecase
	RecaptchaUsecase recaptchaU.RecaptchaUsecase
}

type Services struct {
	TemplateService service.TemplateService
	UrlService      urlS.UrlService
}

type Server struct {
	Mongo        *mongoDb.MongoClient
	Redis        redisDb.Redis
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
	urlRepository := mongoRepository.NewMongoUrlRepository(mongoClient.GetDatabase())

	// usecases
	hashUsecase := hashU.NewHashUsecase(node)
	urlUsecase := urlU.NewUrlUsecase(urlRepository, redisClient)
	rateLimitUsecase := ratelimitU.NewRateLimitUsecase(redisClient)
	recaptchaUsecase := recaptchaU.NewRecaptchaUsecase(recaptchaObj)

	// services
	templateService := service.NewTemplateService(recaptchaKeys.Public)
	urlService := urlS.NewUrlService(urlUsecase, hashUsecase, rateLimitUsecase, recaptchaUsecase)

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
