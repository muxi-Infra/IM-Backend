package redis

import (
	"IM-Backend/cache"
	"IM-Backend/errcode"
	"context"
	"github.com/redis/go-redis/v9"
)

type Reader struct {
	cli *redis.Client
}

func (r *Reader) GetKV(ctx context.Context, kv cache.KV) error {
	res, err := r.cli.Get(ctx, kv.GetStrKey()).Result()
	if err != nil {
		return errcode.ERRGetKV.WrapError(err)
	}
	err = kv.ReadFromStrVal(res)
	if err != nil {
		return errcode.ERRConvertJson.WrapError(err)
	}
	return nil
}

func (r *Reader) GetValFromSet(ctx context.Context, kv cache.KV) ([]string, error) {
	cli, err := r.cli.SMembers(ctx, kv.GetSetKey()).Result()
	if err != nil {
		return nil, errcode.ERRGetSet.WrapError(err)
	}
	return cli, nil
}
