package controller

import (
	"IM-Backend/model"
	"IM-Backend/model/table"
	"context"
	"time"
)

type PostService interface {
	Create(ctx context.Context, svc string, postInfo table.PostInfo) error
	Update(ctx context.Context, svc, userID string, postID uint64, updates map[string]interface{}) error
	GetInfo(ctx context.Context, svc string, postID uint64) (model.PostInfo, error)
	GetLike(ctx context.Context, svc string, postID uint64) (int64, error)
	Delete(ctx context.Context, svc string, userID string, postID uint64) error
	Like(ctx context.Context, svc string, postID uint64, userID string) error
}

type PostIDGenerator interface {
	GeneratePostID(ctx context.Context, svc string) (uint64, error)
}

type CommentService interface {
	Publish(ctx context.Context, svc string, comment table.PostCommentInfo) error
	Update(ctx context.Context, svc, userID string, commentID uint64, updates map[string]interface{}) error
	FindComment(ctx context.Context, svc string, fatherID uint64, cursor time.Time, limit uint) ([]model.PostComment, error)
	GetLike(ctx context.Context, svc string, commentID ...uint64) (map[uint64]int64, error)
	GetChildCommentCnt(ctx context.Context, svc string, commentID ...uint64) (map[uint64]int, error)
	Delete(ctx context.Context, svc string, userID string, commentID uint64) error
	GetCommentUserIDByID(ctx context.Context, svc string, commentID uint64) (string, error)
	Like(ctx context.Context, svc string, postID uint64, commentID uint64, userID string) error
}
type CommentIDGenerator interface {
	GenerateCommentID(ctx context.Context, svc string) (uint64, error)
}
