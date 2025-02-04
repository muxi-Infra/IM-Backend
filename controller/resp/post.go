package resp

import (
	"IM-Backend/model"
	"time"
)

type Post struct {
	ID        uint64                 `json:"id"`
	Content   string                 `json:"content"`
	Author    string                 `json:"author"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Svc       string                 `json:"svc"`
	Extra     map[string]interface{} `json:"extra"`
}

func NewPostResp(post model.PostInfo) Post {
	return Post{
		ID:        post.ID,
		Content:   post.Content,
		Author:    post.Author,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		Svc:       post.Svc,
		Extra:     post.Extra,
	}
}
