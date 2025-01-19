package dao

import (
	"IM-Backend/model/table"
	"context"
	"gorm.io/gorm"
)

type Table interface {
	TableName(svc string) string
	PgCreate(db *gorm.DB, svc string) error
}
type TableTooler interface {
	NewTable(db *gorm.DB, t Table, svc string) error
	CheckTableExist(db *gorm.DB, t Table, svc string) bool
}

type GormWriter interface {
	GetGormTx(ctx context.Context) (tx *gorm.DB)
	Create(ctx context.Context, svc string, t Table) error
	Update(ctx context.Context, svc string, t Table) error
	Delete(ctx context.Context, svc string, t Table) error
	InTx(ctx context.Context, f func(ctx context.Context) error) error
}
type GormReader interface {
	GetPostInfos(ctx context.Context, svc string, ids ...uint64) ([]table.PostInfo, error)
	GetPostLike(ctx context.Context, svc string, id uint64) (int64, error)

	GetPostInfosByExtra(ctx context.Context, svc string, key string, value interface{}) ([]table.PostInfo, error)

	GetCommentIds(ctx context.Context, svc string, id uint64) ([]uint64, error)
	GetCommentInfos(ctx context.Context, svc string, ids ...uint64) ([]table.PostCommentInfo, error)
	GetCommentLike(ctx context.Context, svc string, ids ...uint64) (map[uint64]int64, error)
	GetCommentsByExtra(ctx context.Context, svc string, key string, value interface{}) ([]table.PostCommentInfo, error)
}
