package redis

import (
	"IM-Backend/model"
	"context"
	"testing"
)

func TestReader_GetKV(t *testing.T) {
	cli := initRedis()
	r := &Reader{cli}
	tmp := &model.PostInfo{
		Svc:    "testsvc",
		ID:     1,
		Author: "hello",
	}
	err := r.GetKV(context.Background(), tmp)
	if err != nil {
		t.Error(err)
	}
	t.Log(tmp)
}

func TestReader_GetValFromSet(t *testing.T) {
	cli := initRedis()
	r := &Reader{cli}
	tmp := &model.PostInfo{
		Svc:    "testsvc",
		ID:     1,
		Author: "hello",
	}
	res, err := r.GetValFromSet(context.Background(), tmp)
	if err != nil {
		t.Error(err)
	}
	t.Log(res)
}
