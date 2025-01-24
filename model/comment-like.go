package model

import (
	"encoding/json"
	"fmt"
)

type PostCommentLike struct {
	Svc       string
	CommentID uint64
	Like      int64
}

func (p *PostCommentLike) ReadFromStrVal(jsonStr string) error {
	return json.Unmarshal([]byte(jsonStr), p)
}

func (p *PostCommentLike) GetStrKey() string {
	return fmt.Sprintf("svc:%s:comment:%d:like", p.Svc, p.CommentID)
}

func (p *PostCommentLike) GetStrVal() string {
	val, _ := json.Marshal(p)
	return string(val)
}

func (p *PostCommentLike) GetSetKey() string {
	return "mistake"
}

func (p *PostCommentLike) GetSetVal() string {
	return "mistake"
}
