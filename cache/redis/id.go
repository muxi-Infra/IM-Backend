package redis

import (
	"IM-Backend/errcode"
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type IDer struct {
	client *redis.Client
}

func NewIDer(client *redis.Client) *IDer {
	return &IDer{client: client}
}

func (I *IDer) GeneratePostID(ctx context.Context, svc string) (uint64, error) {
	res, err := I.client.Incr(ctx, I.getPostIDKey(svc)).Result()
	if err != nil {
		return 0, errcode.ERRGenerateID.WrapError(err)
	}
	// 使用秒级时间戳扩大时间范围
	secondsTimeStamp := time.Now().Unix()
	// 组合时间戳（高32位）和自增值（低32位），确保自增值不溢出
	postID := (uint64(secondsTimeStamp) << 32) | (uint64(res) & 0xFFFFFFFF)
	return postID, nil
}

func (I *IDer) GenerateCommentID(ctx context.Context, svc string) (uint64, error) {
	res, err := I.client.Incr(ctx, I.getCommentIDKey(svc)).Result()
	if err != nil {
		return 0, errcode.ERRGenerateID.WrapError(err)
	}
	// 使用秒级时间戳扩大时间范围
	secondsTimeStamp := time.Now().Unix()
	// 组合时间戳（高32位）和自增值（低32位），确保自增值不溢出
	commentID := (uint64(secondsTimeStamp) << 32) | (uint64(res) & 0xFFFFFFFF)
	return commentID, nil
}

func (I *IDer) getPostIDKey(svc string) string {
	return fmt.Sprintf("post_id:%s", svc)
}

func (I *IDer) getCommentIDKey(svc string) string {
	return fmt.Sprintf("comment_id:%s", svc)
}
