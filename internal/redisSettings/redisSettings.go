package redisSettings

import (
	"github.com/go-redis/redis/v8"
	"time"
)

func InitializeRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:            "localhost:6379",
		Password:        "12345",
		DB:              0,               // use default DB
		MaxRetries:      5,               // попыток конекта
		MinRetryBackoff: 2 * time.Second, // время между конектами
	})
}
