package redis

import (
	"context"
	"time"
	"github.com/go-redis/redis/v8"
	"redirect-service/config"
)

var (
	rdb  *redis.Client
	ctx  = context.Background()
)

func InitializeRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisAddress, 
		Password: "",                  
		DB:       0,                   
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}

func SetURLInCache(id string, url string, expiration time.Duration) {
	rdb.Set(ctx, id, url, expiration)
}

func GetURLFromCache(id string) (string, error) {
	return rdb.Get(ctx, id).Result()
}
