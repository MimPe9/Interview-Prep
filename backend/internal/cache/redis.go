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
	log.Println("Initializing Redis connection")

	redisURL := config.GetRedisUrl()
	var options *redis.Options

	if redisURL != "" {
		log.Printf("Using Redis URL: %s", redisURL)
		opts, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Printf("Failed to parse Redis URL: %v", err)
			return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
		}
		options = opts
	} else {
		host := config.GetRedisHost()
		port := config.GetRedisPort()
		log.Printf("Using Redis host: %s, port: %s", host, port)
		options = &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: "", // no password set
			DB:       0,  // use default DB
		}
	}

	client := redis.NewClient(options)
	ctx := context.Background()

	// Test connection
	log.Println("Testing Redis connection...")
	_, err := client.Ping(ctx).Result()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Successfully connected to Redis")
	return &RedisCache{
		client: client,
		ctx:    ctx,
	}, nil
}

func (r *RedisCache) Set(key string, value interface{}, expiration time.Duration) error {
	log.Printf("Setting cache key: %s", key)

	jsonData, err := json.Marshal(value)
	if err != nil {
		log.Printf("Failed to marshal value for key %s: %v", key, err)
		return err
	}

	err = r.client.Set(r.ctx, key, jsonData, expiration).Err()
	if err != nil {
		log.Printf("Failed to set cache key %s: %v", key, err)
		return err
	}

	log.Printf("Cache set successfully for key: %s", key)
	return nil
}

func (r *RedisCache) Get(key string, value interface{}) error {
	log.Printf("Getting cache key: %s", key)

	data, err := r.client.Get(r.ctx, key).Bytes()
	if err == redis.Nil {
		log.Printf("Cache miss for key: %s", key)
		return fmt.Errorf("cache miss")
	} else if err != nil {
		log.Printf("Error getting cache key %s: %v", key, err)
		return err
	}

	err = json.Unmarshal(data, value)
	if err != nil {
		log.Printf("Failed to unmarshal cache data for key %s: %v", key, err)
		return err
	}

	log.Printf("Cache hit for key: %s", key)
	return nil
}

func (r *RedisCache) Delete(key string) error {
	log.Printf("Deleting cache key: %s", key)

	err := r.client.Del(r.ctx, key).Err()
	if err != nil {
		log.Printf("Failed to delete cache key %s: %v", key, err)
		return err
	}

	log.Printf("Cache deleted successfully for key: %s", key)
	return nil
}

func (r *RedisCache) DeletePattern(pattern string) error {
	log.Printf("Deleting cache pattern: %s", pattern)

	keys, err := r.client.Keys(r.ctx, pattern).Result()
	if err != nil {
		log.Printf("Failed to get keys for pattern %s: %v", pattern, err)
		return err
	}

	if len(keys) > 0 {
		log.Printf("Deleting %d keys matching pattern: %s", len(keys), pattern)
		err = r.client.Del(r.ctx, keys...).Err()
		if err != nil {
			log.Printf("Failed to delete keys for pattern %s: %v", pattern, err)
			return err
		}
	} else {
		log.Printf("No keys found for pattern: %s", pattern)
	}

	return nil
}

func (r *RedisCache) Close() error {
	log.Println("Closing Redis connection")
	return r.client.Close()
}
