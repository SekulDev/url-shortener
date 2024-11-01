package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"url-shortener/internal/domain/entity"
)

type UrlRepository interface {
	Create(url *entity.Url) (primitive.ObjectID, error)
	GetByID(id string) (*entity.Url, error)
	GetByShortUrl(url string) (*entity.Url, error)
}
