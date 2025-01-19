package pg

import (
	"IM-Backend/dao"
	"IM-Backend/model/table"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestPgTable_CheckTableExist(t *testing.T) {
	// PostgreSQL 数据库连接配置
	dsn := "host=localhost user=chen password=12345678 dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	// 连接到 PostgreSQL 数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接到数据库失败")
	}
	var tt dao.Table = &table.PostInfo{}
	var pg PgTable
	exist := pg.CheckTableExist(db, tt, "testsvc")
	t.Log(exist)
	exist = pg.CheckTableExist(db, tt, "testsvc1")
	t.Log(exist)
}
