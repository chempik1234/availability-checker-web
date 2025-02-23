package redis

import (
	"context"
	"github.com/go-redis/redis/v7"
	"log"
)

// NewRedisClient try to connect to Redis and get the client
func NewRedisClient(ctx context.Context, redisURL string) (*redis.Client, error) {
	option, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(option)
	err = client.Ping().Err()
	if err != nil {
		return nil, err
	}
	log.Printf("connected to a redis database: %s\n", redisURL)
	return client, nil
}
