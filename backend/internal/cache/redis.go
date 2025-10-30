package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"interviewPrep/backend/config"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisCache() (*RedisCache, error) {
	redisURL := config.GetRedisUrl()
	var options *redis.Options

	if redisURL != "" {
		opts, err := redis.ParseURL(redisURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
		}
		options = opts
	} else {
		options = &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", config.GetRedisHost(), config.GetRedisPort()),
			Password: "", // no password set
			DB:       0,  // use default DB
		}
	}

	client := redis.NewClient(options)
	ctx := context.Background()

	//Test connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, key, jsonData, expiration).Err()
}

func (r *RedisCache) Get(key string, value interface{}) error {
	data, err := r.client.Get(r.ctx, key).Bytes()
	if err == redis.Nil {
		return fmt.Errorf("cache miss")
	} else if err != nil {
		return nil
	}

	return json.Unmarshal(data, value)
}

func (r *RedisCache) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

func (r *RedisCache) DeletePattern(pattern string) error {
	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		return err
	}

	if len(keys) > 0 {
		return r.client.Del(r.ctx, keys...).Err()
	}

	return nil
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
