package redis

import (
	"log"

	redis "github.com/go-redis/redis/v7"
)

// NewClient creates a new redis client
func NewClient(logger *log.Logger) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return client
}
