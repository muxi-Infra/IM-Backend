package pg

import (
	"IM-Backend/dao"
	"IM-Backend/model/table"
	"gorm.io/gorm"
	"sync"
)

var (
	// 使用 sync.Pool 重用对象，减少GC压力
	// 由于使用了很多临时变量来获取表名，所以构建对象池，可以根据并发量来自动扩缩容
	// 防止内存在并发量高的情况下突增
	PostInfoPool = sync.Pool{
		New: func() interface{} {
			return new(table.PostInfo)
		},
	}
	PostLikePool = sync.Pool{
		New: func() interface{} {
			return new(table.PostLikeInfo)
		},
	}
	CommentInfoPool = sync.Pool{
		New: func() interface{} {
			return new(table.PostCommentInfo)
		},
	}
	CommentLikePool = sync.Pool{
		New: func() interface{} {
			return new(table.CommentLikeInfo)
		},
	}
)

type PgTable struct{}

func NewPgTable() *PgTable {
	return &PgTable{}
}

func (p *PgTable) NewTable(db *gorm.DB, t dao.Table, svc string) error {
	return t.PgCreate(db, svc)
}

func (p *PgTable) CheckTableExist(db *gorm.DB, t dao.Table, svc string) bool {
	tableName := t.TableName(svc)
	var count int64
	// 查询 pg_tables 来检查某个表是否存在
	err := db.Raw(`SELECT COUNT(*) 
	FROM pg_tables 
	WHERE schemaname = 'public' AND tablename = ?`, tableName).Scan(&count).Error
	if err != nil {
		return false
	}
	return count > 0
}
