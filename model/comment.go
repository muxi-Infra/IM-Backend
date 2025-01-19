package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type PostComment struct {
	ID           uint64
	UserID       string
	Content      string
	FatherID     uint64
	FatherUserID string
	PostID       uint64
	Time         time.Time
	Svc          string
	Like         int64
	Extra        map[string]interface{}
}

func (c *PostComment) ReadFromStrVal(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), c)
}

func (c *PostComment) GetStrKey() string {
	return fmt.Sprintf("svc:%s:user:%s:post:%d:comment:%d", c.Svc, c.UserID, c.PostID, c.ID)
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
