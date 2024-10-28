package database

import (
	"context"
	"github.com/go-redis/redis"
	"log"
)

type RedisClient struct {
	Client *redis.Client
	ctx    context.Context
}

func NewRedisClient(opts *redis.Options) *RedisClient {
	log.Println("Connected to Redis!")
	return &RedisClient{
		ctx:    context.Background(),
		Client: redis.NewClient(opts),
	}
}

func (r *RedisClient) Disconnect() error {
	log.Println("Disconnected from Redis!")
	err := r.Client.Close()
	if err != nil {
		log.Fatalf("Could not close Redis client: %v", err)
	}
	return err
}
