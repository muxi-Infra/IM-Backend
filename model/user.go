package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	UserTable = "users"
)

type User struct {
	ID string //用户ID，为了适配各个服务，直接用string了
	//Svc uint64 //所属服务ID
	CreatedAt time.Time //创建时间
}

func (u *User) PgCreate(db *gorm.DB, svc string) error {
	tableName := u.TableName(svc)
	sql := fmt.Sprintf(`CREATE TABLE %s (
    id VARCHAR(255),                                    -- 用户ID，使用VARCHAR来存储用户ID，支持各类格式
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- 创建时间，默认为当前时间
    PRIMARY KEY (id),                                    -- 用户ID作为主键，确保唯一性
    INDEX idx_created_at (created_at)                    -- 为 created_at 创建索引，优化基于时间的查询
	);
	`, tableName)
	return db.Exec(sql).Error
}

func (u *User) TableName(svc string) string {
	return svc + "_" + UserTable
}
