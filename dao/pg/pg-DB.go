package pg

import (
	"IM-Backend/dao"
	"gorm.io/gorm"
)

type PgTable struct{}

func (p PgTable) NewTable(db *gorm.DB, t dao.Table, svc string) error {
	return t.PgCreate(db, svc)
}

func (p PgTable) CheckTableExist(db *gorm.DB, t dao.Table, svc string) bool {
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
