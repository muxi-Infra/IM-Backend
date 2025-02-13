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
	ID     uint64 `gorm:"column:id"`      //评论唯一ID
	UserID string `gorm:"column:user_id"` //用户ID
	//SvcID uint64//所属服务ID
	RootID       uint64    `gorm:"column:root_id"`          //该评论所属的根评论ID,如果本身就是根评论，则为0
	FatherID     uint64    `gorm:"column:father_id"`        //该评论所回复的评论，如果没有所回复的评论，则为0
	TargetUserID *string   `gorm:"column:target_user_id"`   //回复评论的用户ID，如果是根评论则为NULL
	PostID       uint64    `gorm:"column:post_id"`          //所属帖子的ID
	Content      string    `gorm:"column:content"`          // 评论内容
	Extra        JSON      `gorm:"type:jsonb;column:extra"` // 其他补充内容，json格式
	CreatedAt    time.Time `gorm:"column:created_at"`       //创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at"`       //修改时间
}

func (p *PostCommentInfo) PgCreate(db *gorm.DB, svc string) error {
	tableName := p.TableName(svc)
	sql := fmt.Sprintf(`
	CREATE TABLE %s (
    id BIGINT,                                         -- 评论唯一ID
    user_id VARCHAR(255),                               -- 用户ID
	root_id BIGINT DEFAULT 0, 							-- 该评论所属的根评论,如果本身就是根评论，则为0
    father_id BIGINT DEFAULT 0,                        -- 该评论所回复的评论，如果没有所回复的评论，则为0
	target_user_id VARCHAR(255) DEFAULT NULL, --回复评论的用户ID，如果是根评论则为NULL
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
	CREATE  INDEX  idx_pc_root_id ON %s (root_id);
`, tableName, tableName, tableName, tableName)

	return db.Exec(sql).Error
}

func (p *PostCommentInfo) TableName(svc string) string {
	return svc + "_" + PostCommentInfoTable
}
