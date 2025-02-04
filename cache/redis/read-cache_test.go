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

func TestReader_MGetKV(t *testing.T) {
	cli := initRedis()
	r := &Reader{cli}
	tmp1 := model.PostInfo{ID: 1, Svc: "testsvc"}
	tmp2 := model.PostInfo{ID: 2, Svc: "testsvc"}
	tmp3 := model.PostInfo{ID: 3, Svc: "testsvc"}
	res := r.MGetKV(context.Background(), &tmp1, &tmp2, &tmp3)
	t.Log(tmp1, tmp2, tmp3)
	t.Log(res)
}
