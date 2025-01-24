package model

import (
	"encoding/json"
	"fmt"
)

type PostLike struct {
	PostID uint64
	Svc    string
	Like   int64
}

func (p *PostLike) ReadFromStrVal(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), p)
}

func (p *PostLike) GetStrKey() string {
	return fmt.Sprintf("svc:%s:post:%d:like", p.Svc, p.PostID)
}

func (p *PostLike) GetStrVal() string {
	val, _ := json.Marshal(p)
	return string(val)
}

func (p *PostLike) GetSetKey() string {
	return "mistake"
}

func (p *PostLike) GetSetVal() string {
	return "mistake"
}
