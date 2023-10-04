package redis

import (
	"github.com/icloudeng/platform-installer/internal/env"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     env.Config.REDIS_URL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
