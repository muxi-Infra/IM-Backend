package model

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestPostInfo_PgCreate(t *testing.T) {
	var info = &PostInfo{}
	// PostgreSQL 数据库连接配置
	dsn := "host=localhost user=chen password=12345678 dbname=mydb port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	// 连接到 PostgreSQL 数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("连接到数据库失败")
	}
	err = info.PgCreate(db, "testsvc")
	if err != nil {
		t.Error(err)
	}
}
