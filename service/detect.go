package service

import (
	"IM-Backend/service/identity"
	"context"
	"time"
)

type DetectSvc struct {
	pfp        chan identity.PostIdentity          //待寻找的post(pending find post)
	pfc        chan identity.CommentIdentity       //待寻找的comment(pending find comment)
	pdc        chan<- identity.CommentIdentity     //待删除的comment(pending delete comment)
	pdpl       chan<- identity.PostLikeIdentity    //待删除的post like(pending delete post like)
	pdcl       chan<- identity.CommentLikeIdentity //待删除的comment like(pending delete comment like)
	finder     TrashFinder
	svcHandler SvcHandler
}

func NewDetectSvc(pfp chan identity.PostIdentity,
	pfc chan identity.CommentIdentity,
	pdc chan<- identity.CommentIdentity,
	pdpl chan<- identity.PostLikeIdentity,
	pdcl chan<- identity.CommentLikeIdentity,
	finder TrashFinder,
	svcHandler SvcHandler) *DetectSvc {
	return &DetectSvc{
		pfp:        pfp,
		pfc:        pfc,
		pdc:        pdc,
		pdpl:       pdpl,
		pdcl:       pdcl,
		finder:     finder,
		svcHandler: svcHandler,
	}
}

func (ds *DetectSvc) Run(ctx context.Context) {
	//另起协程来找寻多余的postID和commentID
	ds.find(ctx)
	for {
		select {
		case post := <-ds.pfp:
			go ds.delComment(ctx, post)
			go ds.delCommentLikeByPost(ctx, post)
			ds.delPostLike(ctx, post)
		case comment := <-ds.pfc:
			ds.delCommentLikeByComment(ctx, comment)
		case <-ctx.Done():
			return
		}
	}
}

func (ds *DetectSvc) find(ctx context.Context) {
	go func(ctxx context.Context) {
		for {
			select {
			case <-time.After(10 * time.Minute):
				svcs := ds.svcHandler.GetAllServices()
				for _, svc := range svcs {
					pss := ds.findPost(ctx, svc)
					for _, post := range pss {
						ds.pfp <- post
					}
				}
			case <-ctxx.Done():
				return
			}
		}
	}(ctx)
	go func(ctxx context.Context) {
		for {
			select {
			case <-time.Tick(15 * time.Minute):
				svcs := ds.svcHandler.GetAllServices()
				for _, svc := range svcs {
					ccs := ds.findComment(ctx, svc)
					for _, comment := range ccs {
						ds.pfc <- comment
					}
				}
			case <-ctxx.Done():
				return
			}
		}
	}(ctx)
}

func (ds *DetectSvc) findPost(ctx context.Context, svc string) []identity.PostIdentity {
	res := ds.finder.FindTrashPostID(ctx, svc)
	if len(res) == 0 {
		return nil
	}
	var tmp = make([]identity.PostIdentity, len(res))
	for k := range res {
		tmp[k] = identity.PostIdentity{
			Svc:    svc,
			PostID: res[k],
		}
	}
	return tmp
}
func (ds *DetectSvc) findComment(ctx context.Context, svc string) []identity.CommentIdentity {
	res := ds.finder.FindTrashCommentID(ctx, svc)
	if len(res) == 0 {
		return nil
	}
	var tmp = make([]identity.CommentIdentity, len(res))
	for k := range res {
		tmp[k] = identity.CommentIdentity{
			Svc:       svc,
			CommentID: res[k],
		}
	}
	return tmp
}

func (ds *DetectSvc) delComment(ctx context.Context, post identity.PostIdentity) {
	comments := ds.finder.FindTrashCommentIDByPostID(ctx, post.Svc, post.PostID)
	for _, comment := range comments {
		ds.pdc <- identity.CommentIdentity{
			Svc:       post.Svc,
			CommentID: comment,
		}
	}
}
func (ds *DetectSvc) delCommentLikeByPost(ctx context.Context, post identity.PostIdentity) {
	mp := ds.finder.FindTrashCommentLikeByPostID(ctx, post.Svc, post.PostID)
	for k, v := range mp {
		for _, vv := range v {
			ds.pdcl <- identity.CommentLikeIdentity{
				Svc:       post.Svc,
				CommentID: k,
				UserID:    vv,
			}
		}
	}
}
func (ds *DetectSvc) delCommentLikeByComment(ctx context.Context, comment identity.CommentIdentity) {
	res := ds.finder.FindTrashCommentLikeByCommentID(ctx, comment.Svc, comment.CommentID)
	for _, v := range res {
		ds.pdcl <- identity.CommentLikeIdentity{
			Svc:       comment.Svc,
			CommentID: comment.CommentID,
			UserID:    v,
		}
	}
}
func (ds *DetectSvc) delPostLike(ctx context.Context, post identity.PostIdentity) {
	res := ds.finder.FindTrashPostLikeByPostID(ctx, post.Svc, post.PostID)
	for _, v := range res {
		ds.pdpl <- identity.PostLikeIdentity{
			Svc:    post.Svc,
			PostID: post.PostID,
			UserID: v,
		}
	}
}
