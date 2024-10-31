package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"url-shortener/internal/domain/entity"
)

type UrlRepository interface {
	Create(url *entity.Url) (primitive.ObjectID, error)
	GetByID(id string) (*entity.Url, error)
	GetByShortUrl(url string) (*entity.Url, error)
}

type mongoUrlRepository struct {
	collection *mongo.Collection
}

func NewMongoUrlRepository(db *mongo.Database) UrlRepository {
	return &mongoUrlRepository{
		collection: db.Collection("urls"),
	}
}

func (m *mongoUrlRepository) Create(url *entity.Url) (primitive.ObjectID, error) {
	result, err := m.collection.InsertOne(context.Background(), url)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (m *mongoUrlRepository) GetByID(id string) (*entity.Url, error) {
	idd, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.D{{"_id", idd}}

	var result entity.Url
	err = m.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *mongoUrlRepository) GetByShortUrl(url string) (*entity.Url, error) {
	filter := bson.D{{"shortid", url}}

	var result entity.Url
	err := m.collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
