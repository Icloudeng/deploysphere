package redis

import (
	// "context"
	// "fmt"
	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
