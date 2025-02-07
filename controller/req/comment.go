package req

type PublishCommentQuery struct {
	AppKey string `form:"appKey"`
	Svc    string `form:"svc"`
	PostID uint64 `form:"post_id"`
	UserID string `form:"user_id"`
}

type PublishCommentFormData struct {
	Content string  `form:"content"`
	Extra   *string `form:"extra"`
}

type ReplyCommentQuery struct {
	AppKey          string `form:"appKey"`
	Svc             string `form:"svc"`
	PostID          uint64 `form:"post_id"`
	RootCommentID   uint64 `form:"root_comment_id"`
	FatherCommentID uint64 `form:"father_comment_id"`
	UserID          string `form:"user_id"`
}

type ReplyCommentFormData struct {
	Content string  `form:"content"`
	Extra   *string `form:"extra"`
}

type UpdateCommentQuery struct {
	AppKey    string `form:"appKey"`
	Svc       string `form:"svc"`
	CommentID uint64 `form:"comment_id"`
	UserID    string `form:"user_id"`
}
type UpdateCommentFormData struct {
	Content *string `form:"content"`
	Extra   *string `form:"extra"`
}

type DeleteCommentQuery struct {
	AppKey    string `form:"appKey"`
	Svc       string `form:"svc"`
	CommentID uint64 `form:"comment_id"`
	UserID    string `form:"user_id"`
}

type GetCommentInfoQuery struct {
	AppKey        string `form:"appKey"`
	Svc           string `form:"svc"`
	PostID        uint64 `form:"post_id"`
	RootCommentID uint64 `form:"root_comment_id"`
	Cursor        string `form:"cursor"`
	Limit         uint   `form:"limit"`
}

type GetCommentLikeQuery struct {
	AppKey string `form:"appKey"`
	Svc    string `form:"svc"`
}
type GetCommentLikeJSON struct {
	CommentIDs []uint64 `json:"comment_ids"`
}

type LikeCommentQuery struct {
	AppKey    string `form:"appKey"`
	Svc       string `form:"svc"`
	PostID    uint64 `form:"post_id"`
	CommentID uint64 `form:"comment_id"`
	UserID    string `form:"user_id"`
}
