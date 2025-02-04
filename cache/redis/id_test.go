package redis

import (
	"context"
	"testing"
)

func TestIDer_GeneratePostID(t *testing.T) {
	client := initRedis()
	ider := NewIDer(client)
	postID, err := ider.GeneratePostID(context.Background(), "testsvc")
	if err != nil {
		t.Error(err)
	}
	t.Log(postID)
}
func TestIDer_GenerateCommentID(t *testing.T) {
	client := initRedis()
	ider := NewIDer(client)
	commentID, err := ider.GenerateCommentID(context.Background(), "testsvc")
	if err != nil {
		t.Error(err)
	}
	t.Log(commentID)
}
