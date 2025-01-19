package cache

import "context"

type KV interface {
	ReadFromStrVal(jsonStr string) error

	GetStrKey() string
	GetStrVal() string

	GetSetKey() string
	GetSetVal() string
}
type WriteCache interface {
	SetKV(ctx context.Context, svc string, kv ...KV) error
	AddKVToSet(ctx context.Context, svc string, kv ...KV) error
}
type ReadCache interface {
	GetKV(ctx context.Context, key string) (KV, error)
	GetValFromSet(ctx context.Context, key string) ([]string, error)
}
