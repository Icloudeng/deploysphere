package redis

import (
	"smatflow/platform-installer/pkg/env"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func init() {
	Client = redis.NewClient(&redis.Options{
		Addr:     env.EnvConfig.REDIS_URL,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
