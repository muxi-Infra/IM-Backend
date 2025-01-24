package redis

import (
	"IM-Backend/model"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

func initRedis() *redis.Client {
	cli := redis.NewClient(&redis.Options{
		DB:   0,
		Addr: "localhost:16379",
	})
	pong, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(pong)
	return cli
}

func TestWriter_SetKV(t *testing.T) {
	cli := initRedis()
	//cli.Set(context.Background(), "hello", "Nihao", 3*time.Minute)
	w := &Writer{cli}
	err := w.SetKV(context.Background(), 20*time.Minute, &model.PostInfo{
		ID:        1,
		Content:   "hello world",
		Author:    "hello",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Svc:       "testsvc",
		Extra: map[string]interface{}{
			"image": "hello",
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestWriter_AddKVToSet(t *testing.T) {
	cli := initRedis()
	w := &Writer{cli}
	err := w.AddKVToSet(context.Background(), 20*time.Minute, &model.PostInfo{
		Svc:    "testsvc",
		ID:     1,
		Author: "hello",
	})
	if err != nil {
		t.Error(err)
	}
}
