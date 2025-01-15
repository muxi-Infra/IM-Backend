package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	PostInfoTable = "post_info"
)

// 帖子基本信息表
type PostInfo struct {
	ID uint64 //帖子唯一ID
	//SvcID  string//所属服务ID
	Content   string                 //内容
	Author    string                 //作者ID
	Extra     map[string]interface{} // 其他补充内容，json格式
	CreatedAt time.Time              //创建时间
	UpdatedAt time.Time              //修改时间
}

func (pi *PostInfo) PgCreate(db *gorm.DB, svc string) error {
	tableName := pi.TableName(svc)
	sql := fmt.Sprintf(`CREATE TABLE %s (
    id BIGINT PRIMARY KEY,              -- 帖子唯一ID (使用 BIGINT 存储 id，保证可以存储较大的数字)
    content TEXT,                       -- 内容 (TEXT 类型用于存储较长的字符串)
    author VARCHAR(255),                -- 作者ID (这里假设为字符串，长度限制为 255)
    extra JSONB,                        -- 其他补充内容 (使用 JSONB 存储 JSON 格式的数据)
    created_at TIMESTAMP WITH TIME ZONE, -- 创建时间 (使用 TIMESTAMP WITH TIME ZONE 支持时区)
    updated_at TIMESTAMP WITH TIME ZONE  -- 修改时间 (使用 TIMESTAMP WITH TIME ZONE 支持时区)
	);
	`, tableName)

	return db.Exec(sql).Error
}

func (pi *PostInfo) TableName(svc string) string {
	return svc + "_" + PostInfoTable
}
