package table

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	PostCommentInfoTable = "post_comment_info"
)

// 帖子评论基本信息表
type PostCommentInfo struct {
	ID     uint64 //评论唯一ID
	UserID string //用户ID
	//SvcID uint64//所属服务ID
	FatherID     uint64    //该评论的父亲评论，如果没有父亲，则为0
	TargetUserID string    //回复评论的用户ID，如果是根评论则为none
	PostID       uint64    //所属帖子的ID
	Content      string    // 评论内容
	Extra        JSON      `gorm:"type:jsonb;column:extra"` // 其他补充内容，json格式
	CreatedAt    time.Time //创建时间
	UpdatedAt    time.Time //修改时间
}

func (p *PostCommentInfo) PgCreate(db *gorm.DB, svc string) error {
	tableName := p.TableName(svc)
	sql := fmt.Sprintf(`
	CREATE TABLE %s (
    id BIGINT,                                         -- 评论唯一ID
    user_id VARCHAR(255),                               -- 用户ID
    father_id BIGINT DEFAULT 0,                        -- 父评论ID，默认为0表示没有父评论
	targert_user_id VARCHAR(255), --回复评论的用户ID，如果是根评论则为none
    post_id BIGINT,                                    -- 所属帖子的ID
    content TEXT,                                      -- 评论内容
    extra JSONB,                                       -- 其他补充内容，存储为JSON格式
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- 创建时间，默认为当前时间
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- 修改时间，默认为当前时间
    PRIMARY KEY (id)                                   -- 评论唯一ID作为主键
	);
	-- 创建索引
	CREATE  INDEX  idx_pc_post_id ON %s (post_id);
	CREATE  INDEX  idx_pc_user_id  ON %s (user_id);
	CREATE  INDEX  idx_pc_father_id ON %s (father_id);
`, tableName, tableName, tableName, tableName)

	return db.Exec(sql).Error
}

func (p *PostCommentInfo) TableName(svc string) string {
	return svc + "_" + PostCommentInfoTable
}
