package table

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

const (
	LikeTypePost    = "post_like"
	LikeTypeComment = "comment_like"
)

// 帖子的点赞表
type PostLikeInfo struct {
	PostID    uint64    `gorm:"column:post_id"` //帖子ID
	UserID    string    `gorm:"column:user_id"` //用户ID
	CreatedAt time.Time //创建时间
}

func (p *PostLikeInfo) PgCreate(db *gorm.DB, svc string) error {
	tableName := p.TableName(svc)
	sql := fmt.Sprintf(`
	CREATE TABLE %s (
		post_id BIGINT,                                  -- 帖子ID
		user_id VARCHAR(255),                            -- 用户ID
		created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- 创建时间，默认为当前时间
		PRIMARY KEY (post_id, user_id)                   -- 联合主键：确保一个帖子和用户的唯一组合
	);

	-- 创建索引
	CREATE INDEX idx_pli_post_id ON %s (post_id);
	CREATE INDEX idx_pli_user_id ON %s (user_id);
	`, tableName, tableName, tableName)
	return db.Exec(sql).Error
}

func (p *PostLikeInfo) TableName(svc string) string {
	return svc + "_" + LikeTypePost
}

// 评论的点赞表
type CommentLikeInfo struct {
	CommentID uint64    //评论ID
	UserID    string    //用户ID
	CreatedAt time.Time //创建时间
}

func (c *CommentLikeInfo) PgCreate(db *gorm.DB, svc string) error {
	tableName := c.TableName(svc)
	sql := fmt.Sprintf(`
	CREATE TABLE %s (
    comment_id BIGINT,                                  -- 评论ID
    user_id VARCHAR(255),                                -- 用户ID
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- 创建时间，默认为当前时间
    PRIMARY KEY (comment_id, user_id)                  -- 联合主键：确保一个用户只能对每条评论点赞一次
	);
	-- 创建索引
	CREATE INDEX idx_cli_comment_id ON %s (comment_id);
	CREATE INDEX idx_cli_user_id ON %s (user_id);
	`, tableName, tableName, tableName)
	return db.Exec(sql).Error
}

func (c *CommentLikeInfo) TableName(svc string) string {
	return svc + "_" + LikeTypeComment
}
