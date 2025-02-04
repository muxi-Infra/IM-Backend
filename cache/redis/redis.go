package redis

import (
	"IM-Backend/configs"
	"context"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(cf configs.AppConf) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     cf.Cache.Addr,
		Password: cf.Cache.Password,
	})
	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return cli
}
