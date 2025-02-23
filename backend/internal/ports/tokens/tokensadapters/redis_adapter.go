package tokensadapters

import (
	"errors"
	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
)

type TokensRepositoryRedis struct {
	RedisClient *redis.Client
}

func NewTokensRepositoryRedis(redisClient *redis.Client) *TokensRepositoryRedis {
	return &TokensRepositoryRedis{RedisClient: redisClient}
}

func (t TokensRepositoryRedis) Check(token string) (bool, error) {
	commandResult := t.RedisClient.Get(token)
	if commandResult.Err() != nil {
		if errors.Is(commandResult.Err(), redis.Nil) {
			return false, nil
		}
		return false, commandResult.Err()
	}
	return true, nil
}

func (t TokensRepositoryRedis) Create() (string, error) {
	var token string
	for {
		token = uuid.New().String()
		alreadyExists, _ := t.Check(token)
		if !alreadyExists {
			break
		}
	}

	commandResult := t.RedisClient.Set(token, "true", 0)
	if commandResult.Err() != nil {
		return "", commandResult.Err()
	}
	return token, nil
}

func (t TokensRepositoryRedis) Delete(token string) error {
	result := t.RedisClient.Del(token)
	if result.Err() != nil {
		return result.Err()
	}
	return nil
}
