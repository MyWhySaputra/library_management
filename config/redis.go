package config

import (
    "github.com/go-redis/redis/v8"
    "os"
)

func ConnectRedis() *redis.Client {
    rdb := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_ADDR"),
        Password: "", // no password set
        DB:       0,  // use default DB
    })

    return rdb
}