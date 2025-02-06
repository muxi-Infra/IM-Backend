package redis

import (
	"IM-Backend/configs"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

func NewRedisClient(cf configs.AppConf) *redis.Client {
	cli := redis.NewClient(&redis.Options{
		Addr:     cf.Cache.Addr,
		Password: cf.Cache.Password,
	})
	_, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(fmt.Sprintf("connect redis failed: %v", err))
	}
	log.Println("connect to redis successfully")
	return cli
}
