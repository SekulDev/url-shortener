package database

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"os"
	"testing"
)

// TestNewMongoClient_Success tests successful connection to MongoDB
func TestNewMongoClient_Success(t *testing.T) {
	// Replace with actual MongoDB URI for integration testing
	uri := os.Getenv("MONGO_URL")
	databaseName := os.Getenv("MONGO_DATABASE")

	mc, err := NewMongoClient(uri, databaseName)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	defer mc.Disconnect()

	if mc.GetDatabase() == nil {
		t.Error("expected database instance, got nil")
	}
}

// TestNewMongoClient_InvalidURI tests that an invalid URI returns an error
func TestNewMongoClient_InvalidURI(t *testing.T) {
	uri := "invalidURI"
	databaseName := "testdb"

	mc, err := NewMongoClient(uri, databaseName)
	if err == nil {
		t.Fatal("expected an error due to invalid URI, got none")
	}

	if mc != nil {
		t.Error("expected MongoClient instance to be nil on error")
	}
}

// TestMongoClient_Disconnect tests that Disconnect closes the MongoDB client connection
func TestMongoClient_Disconnect(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	defer mt.Client.Disconnect(context.TODO())

	mc := &MongoClient{client: mt.Client, databaseName: "testdb"}

	if err := mc.Disconnect(); err != nil {
		t.Errorf("expected no error on disconnect, got %v", err)
	}
}
