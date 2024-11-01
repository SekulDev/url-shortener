package repository

import (
	"testing"
	"url-shortener/internal/domain/entity"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
)

func TestMongoUrlRepository_Create(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Create a new URL", func(mt *mtest.T) {
		repo := &mongoUrlRepository{collection: mt.Coll}

		urlEntity := &entity.Url{
			LongUrl: "https://example.com",
			ShortId: "abc123",
		}

		mt.AddMockResponses(mtest.CreateSuccessResponse())
		id, err := repo.Create(urlEntity)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if id == primitive.NilObjectID {
			t.Error("expected a valid ObjectID, got NilObjectID")
		}
	})
}

func TestMongoUrlRepository_GetByID(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Get URL by ID", func(mt *mtest.T) {
		repo := &mongoUrlRepository{collection: mt.Coll}

		// Creating a mock ObjectID and URL entity
		mockID := primitive.NewObjectID()
		urlEntity := &entity.Url{
			LongUrl: "https://example.com",
			ShortId: "abc123",
		}

		// Adding mock response for the FindOne operation
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.urls", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: mockID},
			{Key: "longurl", Value: urlEntity.LongUrl},
			{Key: "shortid", Value: urlEntity.ShortId},
		}))

		// Fetching by ID
		result, err := repo.GetByID(mockID.Hex())
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.LongUrl != urlEntity.LongUrl {
			t.Errorf("expected LongUrl %s, got %s", urlEntity.LongUrl, result.LongUrl)
		}

		if result.ShortId != urlEntity.ShortId {
			t.Errorf("expected ShortId %s, got %s", urlEntity.ShortId, result.ShortId)
		}
	})

}

func TestMongoUrlRepository_GetByShortUrl(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))

	mt.Run("Get URL by Short URL", func(mt *mtest.T) {
		repo := &mongoUrlRepository{collection: mt.Coll}

		// Creating a mock ObjectID and URL entity
		mockID := primitive.NewObjectID()
		urlEntity := &entity.Url{
			LongUrl: "https://example.com",
			ShortId: "abc123",
		}

		// Adding mock response for the FindOne operation by short URL
		mt.AddMockResponses(mtest.CreateCursorResponse(1, "test.urls", mtest.FirstBatch, bson.D{
			{Key: "_id", Value: mockID},
			{Key: "longurl", Value: urlEntity.LongUrl},
			{Key: "shortid", Value: urlEntity.ShortId},
		}))

		// Fetching by short URL
		result, err := repo.GetByShortUrl(urlEntity.ShortId)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if result.ShortId != urlEntity.ShortId {
			t.Errorf("expected ShortId %s, got %s", urlEntity.ShortId, result.ShortId)
		}
		if result.LongUrl != urlEntity.LongUrl {
			t.Errorf("expected LongUrl %s, got %s", urlEntity.LongUrl, result.LongUrl)
		}
	})

}
