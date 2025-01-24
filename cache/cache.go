package cache

type KV interface {
	ReadFromStrVal(jsonStr string) error

	GetStrKey() string
	GetStrVal() string

	GetSetKey() string
	GetSetVal() string
}
