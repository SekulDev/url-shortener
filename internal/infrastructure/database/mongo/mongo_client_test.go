package database

import (
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
