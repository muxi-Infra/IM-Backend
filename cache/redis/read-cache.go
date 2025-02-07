package redis

import (
	"IM-Backend/cache"
	"IM-Backend/errcode"
	"IM-Backend/global"
	"context"
	"github.com/redis/go-redis/v9"
)

type Reader struct {
	cli *redis.Client
}

func NewReader(cli *redis.Client) *Reader {
	return &Reader{cli: cli}
}

func (r *Reader) MGetKV(ctx context.Context, kv ...cache.KV) []bool {
	if len(kv) == 0 {
		return nil
	}
	keys := make([]string, 0, len(kv))
	oks := make([]bool, 0, len(kv))
	for _, k := range kv {
		keys = append(keys, k.GetStrKey())
	}
	res, err := r.cli.MGet(ctx, keys...).Result()
	if err != nil {
		return nil
	}
	for k, v := range kv {
		valStr, ok := res[k].(string)
		if !ok {
			oks = append(oks, false)
			continue
		}
		err := v.ReadFromStrVal(valStr)
		if err != nil {
			oks = append(oks, false)
			continue
		}
		oks = append(oks, true)
	}
	return oks
}

func (r *Reader) GetKV(ctx context.Context, kv cache.KV) error {
	res, err := r.cli.Get(ctx, kv.GetStrKey()).Result()
	if err != nil {
		global.Log.Warnf("get key:%v from cache failed: %v", kv.GetStrKey(), err)
		return errcode.ERRGetKV.WrapError(err)
	}
	err = kv.ReadFromStrVal(res)
	if err != nil {
		global.Log.Warnf("read from val[%v] in cache failed: %v", res, err)
		return errcode.ERRConvertJson.WrapError(err)
	}
	return nil
}

func (r *Reader) GetValFromSet(ctx context.Context, kv cache.KV) ([]string, error) {
	cli, err := r.cli.SMembers(ctx, kv.GetSetKey()).Result()
	if err != nil {
		global.Log.Warnf("get set[key:%v] from cache failed: %v", kv.GetSetKey(), err)
		return nil, errcode.ERRGetSet.WrapError(err)
	}
	return cli, nil
}
