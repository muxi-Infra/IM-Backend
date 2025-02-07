package resp

import (
	"IM-Backend/model"
	"IM-Backend/pkg"
)

type Post struct {
	ID      uint64                 `json:"id"`
	Title   string                 `json:"title"`
	Content string                 `json:"content"`
	Author  string                 `json:"author"`
	Time    string                 `json:"time"`
	Svc     string                 `json:"svc"`
	Extra   map[string]interface{} `json:"extra"`
}

func NewPostResp(post model.PostInfo) Post {
	return Post{
		ID:      post.ID,
		Title:   post.Title,
		Content: post.Content,
		Author:  post.Author,
		Time:    pkg.FormatTimeInShanghai(post.CreatedAt),
		Svc:     post.Svc,
		Extra:   post.Extra,
	}
}
