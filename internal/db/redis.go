package db

import (
	"log"

	"github.com/redis/go-redis/v9"
	"shirinec.com/config"
)

var Redis *redis.Client

func NewRedis() {
	opts, err := redis.ParseURL(config.AppConfig.RedisURL)
    if err != nil {
        log.Panic(err)
    }
    Redis = redis.NewClient(opts)
}
