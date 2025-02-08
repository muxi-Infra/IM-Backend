package req

type PublishPostQuery struct {
	AppKey string `form:"appKey"`
	Svc    string `form:"svc"`
}
type PublishPostFormData struct {
	Author  string  `form:"author"`
	Title   string  `form:"title"`
	Content string  `form:"content"`
	Extra   *string `form:"extra"`
}

type GetPostInfoQuery struct {
	AppKey string `form:"appKey"`
	Svc    string `form:"svc"`
	PostID uint64 `form:"post_id"`
}

type DeletePostQuery struct {
	AppKey string `form:"appKey"`
	Svc    string `form:"svc"`
	PostID uint64 `form:"post_id"`
	UserID string `form:"user_id"`
}

type GetPostLikeQuery struct {
	AppKey string `form:"appKey"`
	Svc    string `form:"svc"`
	PostID uint64 `form:"post_id"`
}

type LikePostQuery struct {
	AppKey string `form:"appKey"`
	Svc    string `form:"svc"`
	PostID uint64 `form:"post_id"`
	UserID string `form:"user_id"`
}

type LikePostJson struct {
	Like bool `json:"like"`
}

type UpdatePostQuery struct {
	AppKey string `form:"appKey"`
	Svc    string `form:"svc"`
	PostID uint64 `form:"post_id"`
	UserID string `form:"user_id"`
}
type UpdatePostFormData struct {
	Content *string `form:"content"`
	Extra   *string `form:"extra"`
}
