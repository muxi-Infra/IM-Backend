package model

import (
	"IM-Backend/model/table"
	"encoding/json"
	"fmt"
	"time"
)

type PostComment struct {
	ID           uint64
	UserID       string
	Content      string
	RootID       uint64
	FatherID     uint64
	TargetUserID *string
	PostID       uint64
	CreatedAt    time.Time //创建时间
	UpdatedAt    time.Time
	Svc          string
	Extra        map[string]interface{}
}

func NewPostComment(t table.PostCommentInfo, svc string) PostComment {
	return PostComment{
		ID:           t.ID,
		UserID:       t.UserID,
		Content:      t.Content,
		RootID:       t.RootID,
		FatherID:     t.FatherID,
		TargetUserID: t.TargetUserID,
		PostID:       t.PostID,
		CreatedAt:    t.CreatedAt,
		UpdatedAt:    t.UpdatedAt,
		Extra:        t.Extra,
		Svc:          svc,
	}
}
func (c *PostComment) ToTable() table.PostCommentInfo {
	return table.PostCommentInfo{
		ID:        c.ID,
		UserID:    c.UserID,
		Content:   c.Content,
		RootID:    c.RootID,
		FatherID:  c.FatherID,
		PostID:    c.PostID,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		Extra:     c.Extra,
	}
}

func (c *PostComment) ReadFromStrVal(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), c)
}

func (c *PostComment) GetStrKey() string {
	return fmt.Sprintf("svc:%s:comment:%d:info", c.Svc, c.ID)
}

func (c *PostComment) GetStrVal() string {
	val, _ := json.Marshal(c)
	return string(val)
}

func (c *PostComment) GetSetKey() string {
	return fmt.Sprintf("svc:%s:user:%s:post:%d:comments", c.Svc, c.UserID, c.PostID)
}

func (c *PostComment) GetSetVal() string {
	return fmt.Sprintf("%d", c.ID)
}
