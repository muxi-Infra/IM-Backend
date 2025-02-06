package model

import (
	"IM-Backend/model/table"
	"encoding/json"
	"fmt"
	"time"
)

type PostInfo struct {
	ID        uint64
	Content   string
	Title     string
	Author    string
	CreatedAt time.Time
	UpdatedAt time.Time
	Svc       string
	Extra     map[string]interface{}
}

func NewPostInfo(pi table.PostInfo, svc string) PostInfo {
	return PostInfo{
		ID:        pi.ID,
		Title:     pi.Title,
		Content:   pi.Content,
		Author:    pi.Author,
		Extra:     pi.Extra,
		CreatedAt: pi.CreatedAt,
		UpdatedAt: pi.UpdatedAt,
		Svc:       svc,
	}
}

func (p *PostInfo) ToTable() table.PostInfo {
	return table.PostInfo{
		ID:        p.ID,
		Title:     p.Title,
		Content:   p.Content,
		Author:    p.Author,
		Extra:     p.Extra,
		CreatedAt: p.CreatedAt,
		UpdatedAt: p.UpdatedAt,
	}
}

func (p *PostInfo) ReadFromStrVal(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), p)
}

func (p *PostInfo) GetStrKey() string {
	return fmt.Sprintf("svc:%s:post:%d", p.Svc, p.ID)
}

func (p *PostInfo) GetStrVal() string {
	val, _ := json.Marshal(p)
	return string(val)
}

func (p *PostInfo) GetSetKey() string {
	return fmt.Sprintf("svc:%s:user:%s:posts", p.Svc, p.Author)
}

func (p *PostInfo) GetSetVal() string {
	return fmt.Sprintf("%d", p.ID)
}
