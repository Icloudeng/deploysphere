package redis

import (
	// "context"
	// "fmt"
	"smatflow/platform-installer/pkg/env"

	"github.com/redis/go-redis/v9"
)

var Client *redis.Client

func init() {
	opts, err := redis.ParseURL(env.EnvConfig.REDIS_URL)
	if err != nil {
		panic(err)
	}

	Client = redis.NewClient(opts)
}
