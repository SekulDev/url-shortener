package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type MongoClient struct {
	client       *mongo.Client
	databaseName string
}

func NewMongoClient(uri string, databaseName string) (*MongoClient, error) {
	clientOptions := options.Client().ApplyURI(uri).SetConnectTimeout(10 * time.Second)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	log.Println("Connected to MongoDB!")
	return &MongoClient{client: client, databaseName: databaseName}, nil
}

func (mc *MongoClient) GetDatabase() *mongo.Database {
	return mc.client.Database(mc.databaseName)
}

func (mc *MongoClient) Disconnect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return mc.client.Disconnect(ctx)
}
