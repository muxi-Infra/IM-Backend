package service

import (
	"IM-Backend/cache"
	"IM-Backend/dao"
	"IM-Backend/model/table"
	"context"
	"gorm.io/gorm"
	"time"
)

//go:generate mockgen -destination=./mocks/mock_svc.go -package=mocks -source=service.go

type GormWriter interface {
	GetGormTx(ctx context.Context) (tx *gorm.DB)
	Create(ctx context.Context, svc string, t dao.Table) error
	Update(ctx context.Context, svc string, t dao.Table) error
	Delete(ctx context.Context, svc string, t dao.Table, where map[string]interface{}) error
	InTx(ctx context.Context, f func(ctx context.Context) error) error
}

type GormPostReader interface {
	GetPostInfos(ctx context.Context, svc string, ids ...uint64) ([]table.PostInfo, error)
	GetPostLike(ctx context.Context, svc string, id uint64) (int64, error)
	GetPostInfosByExtra(ctx context.Context, svc string, key string, value interface{}) ([]table.PostInfo, error)
	CheckPostExist(ctx context.Context, svc string, id uint64) bool
	GetPostCommentIds(ctx context.Context, svc string, id uint64) ([]uint64, error)
}

type GormCommentReader interface {
	GetCommentInfosByID(ctx context.Context, svc string, ids ...uint64) ([]table.PostCommentInfo, error)
	GetCommentLike(ctx context.Context, svc string, ids ...uint64) (map[uint64]int64, error)
	GetChildCommentIDAfterCursor(ctx context.Context, svc string, fatherID uint64, cursor time.Time, limit uint) ([]uint64, error)
	GetChildCommentCnt(ctx context.Context, svc string, commentID ...uint64) (map[uint64]int, error)
	GetUserIDByCommentID(ctx context.Context, svc string, commentID uint64) (string, error)
	CheckCommentExist(ctx context.Context, svc string, commentID ...uint64) map[uint64]bool
}
type WriteCache interface {
	SetKV(ctx context.Context, expire time.Duration, kv ...cache.KV) error
	DelKV(ctx context.Context, kv cache.KV) error
	AddKVToSet(ctx context.Context, expire time.Duration, kv ...cache.KV) error
}
type ReadCache interface {
	GetKV(ctx context.Context, kv cache.KV) error
	MGetKV(ctx context.Context, kv ...cache.KV) []bool
	GetValFromSet(ctx context.Context, kv cache.KV) ([]string, error)
}

type TrashFinder interface {
	FindTrashPostID(ctx context.Context, svc string) []uint64
	FindTrashCommentID(ctx context.Context, svc string) []uint64

	FindTrashCommentIDByPostID(ctx context.Context, svc string, postID uint64) []uint64
	FindTrashPostLikeByPostID(ctx context.Context, svc string, postID uint64) []string

	FindTrashCommentLikeByPostID(ctx context.Context, svc string, postID uint64) map[uint64][]string
	FindTrashCommentLikeByCommentID(ctx context.Context, svc string, commentID uint64) []string
}

type SvcHandler interface {
	GetAllServices() []string
	GetSecretByName(name string) string
}

type TrashCleaner interface {
	DeleteComment(ctx context.Context, svc string, commentID ...uint64) error
	DeleteCommentLike(ctx context.Context, svc string, commentID uint64, userIDs ...string) error
	DeletePostLike(ctx context.Context, svc string, postID uint64, userIDs ...string) error
}
