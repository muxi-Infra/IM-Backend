package pg

import (
	"IM-Backend/dao"
	"IM-Backend/errcode"
	"IM-Backend/model/table"
	"context"
	"gorm.io/gorm"
)

type TrashCleaner struct {
	db *gorm.DB
	tt dao.TableTooler
}

func NewTrashCleaner(db *gorm.DB, tt dao.TableTooler) *TrashCleaner {
	return &TrashCleaner{
		db: db,
		tt: tt,
	}
}

func (t *TrashCleaner) DeleteComment(ctx context.Context, svc string, commentID ...uint64) error {
	if len(commentID) == 0 {
		return nil
	}
	tmp := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmp)
	if !t.tt.CheckTableExist(t.db, tmp, svc) {
		return errcode.ERRNoTable
	}
	err := t.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("id IN ?", commentID).Delete(&table.PostCommentInfo{}).Error
	if err != nil {
		return errcode.ERRDeleteData.WrapError(err)
	}
	return nil
}

func (t *TrashCleaner) DeleteCommentLike(ctx context.Context, svc string, commentID uint64, userIDs ...string) error {
	if len(userIDs) == 0 {
		return nil
	}
	tmp := CommentLikePool.Get().(*table.CommentLikeInfo)
	defer CommentLikePool.Put(tmp)
	if !t.tt.CheckTableExist(t.db, tmp, svc) {
		return errcode.ERRNoTable
	}
	err := t.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("comment_id =? AND user_id IN?", commentID, userIDs).Delete(&table.CommentLikeInfo{}).Error
	if err != nil {
		return errcode.ERRDeleteData.WrapError(err)
	}
	return nil
}

func (t *TrashCleaner) DeletePostLike(ctx context.Context, svc string, postID uint64, userIDs ...string) error {
	if len(userIDs) == 0 {
		return nil
	}
	tmp := PostLikePool.Get().(*table.PostLikeInfo)
	defer PostLikePool.Put(tmp)
	if !t.tt.CheckTableExist(t.db, tmp, svc) {

		return errcode.ERRNoTable
	}
	err := t.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("post_id =? AND user_id IN?", postID, userIDs).Delete(&table.PostLikeInfo{}).Error

	if err != nil {
		return errcode.ERRDeleteData.WrapError(err)
	}
	return nil
}
