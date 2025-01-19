package model

import (
	"encoding/json"
	"fmt"
	"time"
)

type Post struct {
	ID        uint64
	Content   string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Svc       string
	Like      int64
	Extra     map[string]interface{}
}

func (p *Post) ReadFromStrVal(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), p)
}

func (p *Post) GetStrKey() string {
	return fmt.Sprintf("svc:%s:user:%s:post:%d", p.Svc, p.Author, p.Like)
}

func (p *Post) GetStrVal() string {
	val, _ := json.Marshal(p)
	return string(val)
}

func (p *Post) GetSetKey() string {
	return fmt.Sprintf("svc:%s:user:%s:posts", p.Svc, p.Author)
}

func (p *Post) GetSetVal() string {
	return fmt.Sprintf("%d", p.ID)
}
