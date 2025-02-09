package resp

import (
	"IM-Backend/model"
	"IM-Backend/pkg"
)

type Comment struct {
	ID           uint64                 `json:"id"`
	UserID       string                 `json:"user_id"`
	Content      string                 `json:"content"`
	RootID       uint64                 `json:"root_id"`
	FatherID     uint64                 `json:"father_id"`
	TargetUserID *string                `json:"target_user_id"`
	PostID       uint64                 `json:"post_id"`
	Time         string                 `json:"time"`
	Svc          string                 `json:"svc"`
	ChildNum     *int                   `json:"child_num"`
	Extra        map[string]interface{} `json:"extra"`
}

func NewCommentResp(c model.PostComment, childNum *int) Comment {
	return Comment{
		ID:           c.ID,
		UserID:       c.UserID,
		Content:      c.Content,
		RootID:       c.RootID,
		FatherID:     c.FatherID,
		TargetUserID: c.TargetUserID,
		PostID:       c.PostID,
		Time:         pkg.FormatTimeInShanghai(c.CreatedAt),
		Svc:          c.Svc,
		ChildNum:     childNum,
		Extra:        c.Extra,
	}
}
