package database

import (
	"context"
	"github.com/go-redis/redis"
	"log"
	"time"
)

type Redis interface {
	Close() error
	Get(key string) (string, error)
	Set(key string, value interface{}, expiration time.Duration) error
	Exists(keys ...string) (int64, error)
}

type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(key).Result()
}

func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(key, value, expiration).Err()
}

func (r *RedisClient) Exists(keys ...string) (int64, error) {
	return r.client.Exists(keys...).Result()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func NewRedisClient(opts *redis.Options) Redis {
	log.Println("Connected to Redis!")
	return &RedisClient{
		ctx:    context.Background(),
		client: redis.NewClient(opts),
	}
}

func (r *RedisClient) Disconnect() error {
	log.Println("Disconnected from Redis!")
	err := r.client.Close()
	if err != nil {
		log.Fatalf("Could not close Redis client: %v", err)
	}
	return err
}
