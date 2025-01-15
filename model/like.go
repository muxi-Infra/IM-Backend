package model

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
	PostID    uint64    //帖子ID
	UserID    string    //用户ID
	CreatedAt time.Time //创建时间
}

func (p *PostLikeInfo) PgCreate(db *gorm.DB, svc string) error {
	tableName := p.TableName(svc)
	sql := fmt.Sprintf(`CREATE TABLE %s (
    post_id BIGINT,                                  -- 帖子ID
    user_id VARCHAR(255),                            -- 用户ID
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- 创建时间，默认为当前时间
    PRIMARY KEY (post_id, user_id),                 -- 联合主键：确保一个帖子和用户的唯一组合
    INDEX idx_post_id (post_id),                    -- 为 post_id 创建索引
    INDEX idx_user_id (user_id)                     -- 为 user_id 创建索引
	);
	`, tableName)
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
	sql := fmt.Sprintf(`CREATE TABLE %s (
    comment_id BIGINT,                                  -- 评论ID
    user_id VARCHAR(255),                                -- 用户ID
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,  -- 创建时间，默认为当前时间
    PRIMARY KEY (comment_id, user_id),                  -- 联合主键：确保一个用户只能对每条评论点赞一次
    INDEX idx_comment_id (comment_id),                   -- 为 comment_id 创建索引
    INDEX idx_user_id (user_id)                          -- 为 user_id 创建索引
	);`, tableName)
	return db.Exec(sql).Error
}

func (c *CommentLikeInfo) TableName(svc string) string {
	return svc + "_" + LikeTypeComment
}
