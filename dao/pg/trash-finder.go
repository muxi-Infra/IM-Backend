package pg

import (
	"IM-Backend/dao"
	"IM-Backend/model/table"
	"IM-Backend/pkg"
	"context"
	"fmt"
	"gorm.io/gorm"
)

type TrashFinder struct {
	db *gorm.DB
	tt dao.TableTooler
}

func NewTrashFinder(db *gorm.DB, tt dao.TableTooler) *TrashFinder {
	return &TrashFinder{db: db, tt: tt}
}

func (t *TrashFinder) FindTrashPostID(ctx context.Context, svc string) []uint64 {

	tmpPostInfo := PostInfoPool.Get().(*table.PostInfo)
	defer PostInfoPool.Put(tmpPostInfo)
	if !t.tt.CheckTableExist(t.db, tmpPostInfo, svc) {
		return nil
	}
	tn := tmpPostInfo.TableName(svc)
	//从postLike中查找垃圾postID(即存在于postLike但不存在于postInfo的)
	res := t.findTrashPostIDJoinPostLike(ctx, svc, tn)

	//从comment中查找垃圾postID
	res = append(res, t.findTrashPostIDJoinComment(ctx, svc, tn)...)

	//从commentLike中查找垃圾postID
	res = append(res, t.findTrashPostIDJoinCommentLike(ctx, svc, tn)...)

	//去重
	return pkg.Unique(res)
}

func (t *TrashFinder) FindTrashCommentID(ctx context.Context, svc string) []uint64 {
	tmpComment := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmpComment)
	if !t.tt.CheckTableExist(t.db, tmpComment, svc) {
		return nil
	}
	tn := tmpComment.TableName(svc)

	//从commentLike中获取垃圾commentID
	res := t.findTrashCommentIDJoinCommentLike(ctx, svc, tn)
	return res
}

func (t *TrashFinder) FindTrashCommentIDByPostID(ctx context.Context, svc string, postID uint64) []uint64 {
	tmpComment := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmpComment)
	if !t.tt.CheckTableExist(t.db, tmpComment, svc) {
		return nil
	}
	var res []uint64
	err := t.db.WithContext(ctx).Table(tmpComment.TableName(svc)).
		Select("id").Where("post_id = ?", postID).
		Pluck("id", &res).Error
	if err != nil {
		return nil
	}
	return res
}

func (t *TrashFinder) FindTrashPostLikeByPostID(ctx context.Context, svc string, postID uint64) []string {
	tmpPostLike := PostLikePool.Get().(*table.PostLikeInfo)
	defer PostLikePool.Put(tmpPostLike)
	if !t.tt.CheckTableExist(t.db, tmpPostLike, svc) {
		return nil
	}
	var res []string
	err := t.db.WithContext(ctx).Table(tmpPostLike.TableName(svc)).
		Select("user_id").Where("post_id = ?", postID).
		Pluck("user_id", &res).Error
	if err != nil {
		return nil
	}
	return res
}

func (t *TrashFinder) FindTrashCommentLikeByPostID(ctx context.Context, svc string, postID uint64) map[uint64][]string {
	tmpCommentLike := CommentLikePool.Get().(*table.CommentLikeInfo)
	defer CommentLikePool.Put(tmpCommentLike)
	if !t.tt.CheckTableExist(t.db, tmpCommentLike, svc) {
		return nil
	}
	type Result struct {
		CommentID uint64
		UserID    string
	}
	var res []Result
	err := t.db.WithContext(ctx).Table(tmpCommentLike.TableName(svc)).
		Select("comment_id, user_id").Where("post_id = ?", postID).
		Scan(&res).Error
	if err != nil {
		return nil
	}
	mp := make(map[uint64][]string)
	for _, v := range res {
		mp[v.CommentID] = append(mp[v.CommentID], v.UserID)
	}
	return mp
}

func (t *TrashFinder) FindTrashCommentLikeByCommentID(ctx context.Context, svc string, commentID uint64) []string {
	tmpCommentLike := CommentLikePool.Get().(*table.CommentLikeInfo)
	defer CommentLikePool.Put(tmpCommentLike)
	if !t.tt.CheckTableExist(t.db, tmpCommentLike, svc) {
		return nil
	}
	var res []string
	err := t.db.WithContext(ctx).Table(tmpCommentLike.TableName(svc)).
		Select("user_id").Where("comment_id =?", commentID).
		Pluck("user_id", &res).Error
	if err != nil {
		return nil
	}
	return res
}

func (t *TrashFinder) findTrashPostIDJoinPostLike(ctx context.Context, svc, postInfoTableName string) []uint64 {
	tmpPostLike := PostLikePool.Get().(*table.PostLikeInfo)
	defer PostLikePool.Put(tmpPostLike)

	if !t.tt.CheckTableExist(t.db, tmpPostLike, svc) {
		return nil
	}
	tn := tmpPostLike.TableName(svc)
	var res []uint64
	err := t.db.WithContext(ctx).Table(tn).Select("post_id").
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.post_id = %s.id", postInfoTableName, tn, postInfoTableName)).
		Where("post_id IS NULL").Pluck("post_id", &res).Error
	if err != nil {
		return nil
	}
	return res
}

func (t *TrashFinder) findTrashPostIDJoinComment(ctx context.Context, svc, postInfoTableName string) []uint64 {
	tmpComment := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmpComment)
	if !t.tt.CheckTableExist(t.db, tmpComment, svc) {
		return nil
	}
	tn := tmpComment.TableName(svc)
	var res []uint64
	err := t.db.WithContext(ctx).Table(tn).Select("post_id").
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.post_id = %s.id", postInfoTableName, tn, postInfoTableName)).
		Where("post_id IS NULL").Pluck("post_id", &res).Error
	if err != nil {
		return nil
	}
	return res
}

func (t *TrashFinder) findTrashPostIDJoinCommentLike(ctx context.Context, svc, postInfoTableName string) []uint64 {
	tmpCommentLike := CommentLikePool.Get().(*table.CommentLikeInfo)
	defer CommentLikePool.Put(tmpCommentLike)
	if !t.tt.CheckTableExist(t.db, tmpCommentLike, svc) {
		return nil
	}
	tn := tmpCommentLike.TableName(svc)
	var res []uint64
	err := t.db.WithContext(ctx).Table(tn).Select("post_id").
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.post_id = %s.id", postInfoTableName, tn, postInfoTableName)).
		Where("post_id IS NULL").Pluck("post_id", &res).Error
	if err != nil {
		return nil
	}
	return res
}

func (t *TrashFinder) findTrashCommentIDJoinCommentLike(ctx context.Context, svc, commentTableName string) []uint64 {
	tmpCommentLike := CommentLikePool.Get().(*table.CommentLikeInfo)
	defer CommentLikePool.Put(tmpCommentLike)
	if !t.tt.CheckTableExist(t.db, tmpCommentLike, svc) {
		return nil
	}
	tn := tmpCommentLike.TableName(svc)
	var res []uint64
	err := t.db.WithContext(ctx).Table(tn).Select("comment_id").
		Joins(fmt.Sprintf("LEFT JOIN %s ON %s.comment_id = %s.id", commentTableName, tn, commentTableName)).
		Where("comment_id IS NULL").Pluck("comment_id", &res).Error
	if err != nil {
		return nil
	}
	return res
}
