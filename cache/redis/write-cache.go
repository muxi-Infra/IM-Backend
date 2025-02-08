package redis

import (
	"IM-Backend/cache"
	"IM-Backend/errcode"
	"IM-Backend/global"
	"context"
	"github.com/redis/go-redis/v9"
	"time"
)

type Writer struct {
	cli *redis.Client
}

func NewWriter(client *redis.Client) *Writer {
	return &Writer{cli: client}
}

func (w *Writer) DelKV(ctx context.Context, kv cache.KV) error {
	key := kv.GetStrKey()
	err := w.cli.Del(ctx, key).Err()
	if err != nil {
		global.Log.Warnf("delete key:%v from cache failed: %v", key, err)
		return errcode.ERRDelKV.WrapError(err)
	}
	return nil
}

func (w *Writer) SetKV(ctx context.Context, expire time.Duration, kv ...cache.KV) error {
	if len(kv) == 0 {
		return nil // nothing to do if no key-value pairs provided.
	}

	// 构建 ARGV 参数
	args := make([]interface{}, 0, len(kv)*2+1)
	args = append(args, int(expire.Seconds())) // 设置过期时间，单位秒
	for _, vv := range kv {
		args = append(args, vv.GetStrKey(), vv.GetStrVal())
	}

	script := `
        local expire = tonumber(ARGV[1])  -- 过期时间，单位秒
        local keys = {}  -- 存储所有的键
        for i = 2, #ARGV, 2 do
            local key = ARGV[i]
            local value = ARGV[i + 1]
            redis.call("SET", key, value)  -- 设置键值
            redis.call("EXPIRE", key, expire)  -- 设置过期时间
        end
        return true
    `

	// 执行 Lua 脚本
	err := w.cli.Eval(ctx, script, nil, args...).Err()
	if err != nil {
		global.Log.Warnf("set kv[%+v] in cache failed: %v", kv, err)
		return errcode.ERRSetKV.WrapError(err)
	}
	return nil
}

func (w *Writer) AddKVToSet(ctx context.Context, expire time.Duration, kv ...cache.KV) error {
	if len(kv) == 0 {
		return nil
	}

	// 构建 ARGV 参数
	args := make([]interface{}, 0, len(kv)*2+1)
	args = append(args, int(expire.Seconds())) // 设置过期时间，单位秒
	for _, vv := range kv {
		args = append(args, vv.GetSetKey(), vv.GetSetVal()) // 添加所有的 key-value
	}

	// Lua 脚本
	script := `
        local expire = tonumber(ARGV[1])  -- 过期时间，单位秒
        for i = 2, #ARGV, 2 do
            local key = ARGV[i]
            local value = ARGV[i + 1]
            redis.call("SADD", key, value)  -- 向集合中添加元素
        end
        redis.call("EXPIRE", ARGV[2], expire)  -- 设置集合的过期时间
        return true
    `

	// 执行 Lua 脚本
	err := w.cli.Eval(ctx, script, nil, args...).Err()
	if err != nil {
		global.Log.Warnf("add kv[%+v] to set in cache failed: %v", kv, err)
		return errcode.ERRAddSet.WrapError(err)
	}
	return nil
}
