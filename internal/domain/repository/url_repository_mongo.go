package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"url-shortener/internal/domain/entities"
)

type UrlRepository interface {
	Create(url *entities.Url) (primitive.ObjectID, error)
	GetByID(id string) (*entities.Url, error)
	GetByShortUrl(url string) (*entities.Url, error)
}

type MongoUrlRepository struct {
	collection *mongo.Collection
}

func NewMongoUrlRepository(db *mongo.Database) *MongoUrlRepository {
	return &MongoUrlRepository{
		collection: db.Collection("urls"),
	}
}

func (m *MongoUrlRepository) Create(url *entities.Url) (primitive.ObjectID, error) {
	result, err := m.collection.InsertOne(context.TODO(), url)
	if err != nil {
		return primitive.NilObjectID, err
	}

	return result.InsertedID.(primitive.ObjectID), nil
}

func (m *MongoUrlRepository) GetByID(id string) (*entities.Url, error) {
	idd, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		panic(err)
	}

	filter := bson.D{{"_id", idd}}

	var result entities.Url
	err = m.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (m *MongoUrlRepository) GetByShortUrl(url string) (*entities.Url, error) {
	filter := bson.D{{"shortid", url}}

	var result entities.Url
	err := m.collection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
