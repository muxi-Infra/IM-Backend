package service

import (
	"IM-Backend/configs"
	"IM-Backend/global"
	"IM-Backend/service/identity"
	"context"
	"sync"
	"time"
)

type CleanSvc struct {
	pdc    <-chan identity.CommentIdentity //待删除的comment(pending delete comment)
	batch1 int

	pdpl   <-chan identity.PostLikeIdentity //待删除的post like(pending delete post like)
	batch2 int

	pdcl   <-chan identity.CommentLikeIdentity //待删除的comment like(pending delete comment like)
	batch3 int

	cleaner TrashCleaner
	lock    sync.Mutex
}

func (c *CleanSvc) Callback(conf configs.AppConf) {
	c.batch1 = conf.Clean.CommentBatch
	c.batch2 = conf.Clean.PostLikeBatch
	c.batch3 = conf.Clean.CommentLikeBatch
	if c.batch1 <= 0 {
		c.batch1 = 10
	}
	if c.batch2 <= 0 {
		c.batch2 = 10
	}
	if c.batch3 <= 0 {
		c.batch3 = 10
	}
	global.Log.Infof("cleanSvc's batch1,batch2,batch3 have been set to [%d,%d,%d]", c.batch1, c.batch2, c.batch3)
}

func NewCleanSvc(pdc <-chan identity.CommentIdentity, pdpl <-chan identity.PostLikeIdentity, pdcl <-chan identity.CommentLikeIdentity, cleaner TrashCleaner, cf configs.AppConf) *CleanSvc {
	cs := &CleanSvc{
		pdc:     pdc,
		pdpl:    pdpl,
		pdcl:    pdcl,
		cleaner: cleaner,
		batch1:  cf.Clean.CommentBatch,
		batch2:  cf.Clean.PostLikeBatch,
		batch3:  cf.Clean.CommentLikeBatch,
	}
	if cs.batch1 <= 0 {
		cs.batch1 = 10
	}
	if cs.batch2 <= 0 {
		cs.batch2 = 10
	}
	if cs.batch3 <= 0 {
		cs.batch3 = 10
	}
	global.Log.Infof("cleanSvc's batch1,batch2,batch3 have been set to [%d,%d,%d]", cs.batch1, cs.batch2, cs.batch3)
	return cs
}

func (c *CleanSvc) Run(ctx context.Context) {
	/*
		buf1,buf2,buf3都是缓冲的切片
		会积攒一定的数目再去一次性处理，防止多次IO
		为了防止有些数据长时间得不到清理，设置过期时间来清理
	*/
	buf1 := make([]identity.CommentIdentity, 0, c.batch1)
	buf2 := make([]identity.PostLikeIdentity, 0, c.batch2)
	buf3 := make([]identity.CommentLikeIdentity, 0, c.batch3)
	for {
		select {
		case comment := <-c.pdc:
			buf1 = append(buf1, comment)
			if len(buf1) >= c.batch1 {
				c.lock.Lock()
				c.delComment(ctx, buf1)
				c.lock.Unlock()
				buf1 = buf1[:0]
			}
		case postLike := <-c.pdpl:
			buf2 = append(buf2, postLike)
			if len(buf2) >= c.batch2 {
				c.lock.Lock()
				c.delPostLike(ctx, buf2)
				c.lock.Unlock()
				buf2 = buf2[:0]
			}
		case commentLike := <-c.pdcl:
			buf3 = append(buf3, commentLike)
			if len(buf3) >= c.batch3 {
				c.lock.Lock()
				c.delCommentLike(ctx, buf3)
				c.lock.Unlock()
				buf3 = buf3[:0]
			}
		case <-time.After(15 * time.Minute):
			c.lock.Lock()
			if len(buf1) > 0 {
				c.delComment(ctx, buf1)
			}
			if len(buf2) > 0 {
				c.delPostLike(ctx, buf2)
			}
			if len(buf3) > 0 {
				c.delCommentLike(ctx, buf3)
			}
			c.lock.Unlock()
			buf1, buf2, buf3 = buf1[:0], buf2[:0], buf3[:0]
		case <-ctx.Done():
			return
		}
	}

}
func (c *CleanSvc) delComment(ctx context.Context, comment []identity.CommentIdentity) {
	var mp = make(map[string][]uint64)

	for _, v := range comment {
		mp[v.Svc] = append(mp[v.Svc], v.CommentID)
	}

	for k, v := range mp {
		if len(v) == 0 {
			continue
		}
		_ = c.cleaner.DeleteComment(ctx, k, v...)
	}

	if len(comment) > 0 {
		global.Log.Infof("cleanSvc has cleaned comment: %+v", comment)
	}
}
func (c *CleanSvc) delPostLike(ctx context.Context, postLike []identity.PostLikeIdentity) {
	type tmp struct {
		svc    string
		postID uint64
	}

	var mp = make(map[tmp][]string)

	for _, v := range postLike {
		mp[tmp{v.Svc, v.PostID}] = append(mp[tmp{v.Svc, v.PostID}], v.UserID)
	}

	for k, v := range mp {
		if len(v) == 0 {
			continue
		}
		_ = c.cleaner.DeletePostLike(ctx, k.svc, k.postID, v...)
	}

	if len(postLike) > 0 {
		global.Log.Infof("cleanSvc has cleaned post_like: %+v", postLike)
	}
}
func (c *CleanSvc) delCommentLike(ctx context.Context, commentLike []identity.CommentLikeIdentity) {
	type tmp struct {
		svc       string
		commentID uint64
	}
	var mp = make(map[tmp][]string)

	for _, v := range commentLike {
		mp[tmp{v.Svc, v.CommentID}] = append(mp[tmp{v.Svc, v.CommentID}], v.UserID)
	}

	for k, v := range mp {
		if len(v) == 0 {
			continue
		}
		_ = c.cleaner.DeleteCommentLike(ctx, k.svc, k.commentID, v...)
	}

	if len(commentLike) > 0 {
		global.Log.Infof("cleanSvc has cleaned comment_like: %+v", commentLike)
	}
}
