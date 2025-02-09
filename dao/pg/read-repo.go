package pg

import (
	"IM-Backend/dao"
	"IM-Backend/errcode"
	"IM-Backend/global"
	"IM-Backend/model/table"
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
)

type ReadRepo struct {
	db *gorm.DB
	tt dao.TableTooler
}

func NewReadRepo(db *gorm.DB, tt dao.TableTooler) *ReadRepo {
	return &ReadRepo{
		db: db,
		tt: tt,
	}
}

func (r *ReadRepo) CheckPostExist(ctx context.Context, svc string, id uint64) bool {
	tmp := PostInfoPool.Get().(*table.PostInfo)
	defer PostInfoPool.Put(tmp)
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return false
	}

	var cnt int64
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("id = ?", id).Count(&cnt).Error
	if err != nil {
		global.Log.Infof("someone is trying to check post_info[svc:%v id:%v] exist in db,but failed: %v", svc, id, err)
		return false
	}
	return cnt > 0
}

func (r *ReadRepo) GetPostInfosByExtra(ctx context.Context, svc string, key string, value interface{}) ([]table.PostInfo, error) {
	tmp := PostInfoPool.Get().(*table.PostInfo)
	defer PostInfoPool.Put(tmp)
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}

	var postInfos []table.PostInfo

	// 使用原生 SQL 查询，根据 extra 字段的 key 查找
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).
		Where("extra ->> ? = ?", key, value). // 使用 jsonb 操作符 ->> 提取 key 对应的值并与 value 比较
		Find(&postInfos).Error

	if err != nil {
		global.Log.Errorf("get post_info[svc:%v] by extra[key:%v val:%v] in db failed: %v", svc, key, value, err)
		return nil, errcode.ERRFindData.WrapError(err)
	}

	return postInfos, nil
}

func (r *ReadRepo) GetCommentsByExtra(ctx context.Context, svc string, key string, value interface{}) ([]table.PostCommentInfo, error) {
	tmp := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmp)

	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}

	var commentInfos []table.PostCommentInfo
	// 使用原生 SQL 查询，根据 extra 字段的 key 查找
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).
		Where("extra ->> ? = ?", key, value). // 使用 jsonb 操作符 ->> 提取 key 对应的值并与 value 比较
		Find(&commentInfos).Error
	if err != nil {
		global.Log.Errorf("get comment_info[svc:%v] by extra[key:%v val:%v] in db failed: %v", svc, key, value, err)
		return nil, errcode.ERRFindData.WrapError(err)
	}
	return commentInfos, nil
}

func (r *ReadRepo) GetPostInfos(ctx context.Context, svc string, ids ...uint64) ([]table.PostInfo, error) {
	if len(ids) == 0 {
		return nil, errcode.ERRFindQueryIsEmpty
	}

	tmp := PostInfoPool.Get().(*table.PostInfo)
	defer PostInfoPool.Put(tmp)

	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}

	var postInfos []table.PostInfo
	res := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("id IN ?", ids).Find(&postInfos)
	if res.Error != nil || res.RowsAffected == 0 {
		err := res.Error
		if res.RowsAffected == 0 {
			err = errors.New("get nothing")
		}
		global.Log.Errorf("get postInfos[svc:%v id:%v] in db failed: %v", svc, ids, err)
		return nil, errcode.ERRFindData.WrapError(res.Error)
	}
	return postInfos, nil
}

func (r *ReadRepo) GetPostLike(ctx context.Context, svc string, id uint64) (int64, error) {
	tmp := PostLikePool.Get().(*table.PostLikeInfo)
	defer PostLikePool.Put(tmp)
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return 0, errcode.ERRNoTable
	}

	var cnt int64
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("post_id = ?", id).Count(&cnt).Error
	if err != nil {
		global.Log.Errorf("get post_like[svc:%v id:%v] from db failed: %v", svc, id, err)
		return 0, errcode.ERRCount.WrapError(err)
	}
	return cnt, nil
}

func (r *ReadRepo) GetPostCommentIds(ctx context.Context, svc string, id uint64) ([]uint64, error) {
	tmp := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmp)

	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}

	var commentIds []uint64
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("post_id = ?", id).Pluck("id", &commentIds).Error
	if err != nil {
		global.Log.Errorf("get comment_ids by post_id[%v] in svc[%v] in db failed: %v", id, svc, err)
		return nil, errcode.ERRFindData.WrapError(err)
	}

	return commentIds, nil
}

func (r *ReadRepo) GetCommentInfosByID(ctx context.Context, svc string, ids ...uint64) ([]table.PostCommentInfo, error) {
	if len(ids) == 0 {
		return nil, errcode.ERRFindQueryIsEmpty
	}

	tmp := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmp)
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}

	var commentInfos []table.PostCommentInfo
	res := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("id IN ?", ids).Find(&commentInfos)
	if res.Error != nil || res.RowsAffected == 0 {
		err := res.Error
		if res.RowsAffected == 0 {
			err = errors.New("get nothing")
		}
		global.Log.Errorf("get commentInfos[svc:%v id:%v] in db failed: %v", svc, ids, err)
		return nil, errcode.ERRFindData.WrapError(res.Error)
	}

	return commentInfos, nil
}

func (r *ReadRepo) GetCommentLike(ctx context.Context, svc string, ids ...uint64) (map[uint64]int64, error) {
	if len(ids) == 0 {
		return nil, errcode.ERRFindQueryIsEmpty
	}

	// 定义返回值 map，用于存储每个 comment_id 的点赞数
	result := make(map[uint64]int64, len(ids))

	// 获取动态表名
	tmp := CommentLikePool.Get().(*table.CommentLikeInfo)
	defer CommentLikePool.Put(tmp)

	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}

	// 执行查询，按 comment_id 分组，并统计每个 comment_id 的点赞数
	rows, err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).
		Select("comment_id, COUNT(*) as like_count").
		Where("comment_id IN ?", ids).
		Group("comment_id").
		Rows()
	if err != nil {
		return nil, errcode.ERRCount.WrapError(err)
	}
	defer rows.Close()

	// 遍历查询结果，将每一行数据扫描到 map 中
	for rows.Next() {
		var commentID uint64
		var likeCount int64
		if err := rows.Scan(&commentID, &likeCount); err != nil {
			return nil, errcode.ERRCount.WrapError(err)
		}
		result[commentID] = likeCount
	}

	// 检查是否有错误
	if err := rows.Err(); err != nil {
		global.Log.Errorf("get comment like[svc:%v] by comment_ids[%v] in db failed: %v", svc, ids, err)
		return nil, errcode.ERRCount.WrapError(err)
	}

	return result, nil
}

func (r *ReadRepo) GetChildCommentIDAfterCursor(ctx context.Context, svc string, rootID uint64, cursor time.Time, limit uint) ([]uint64, error) {
	tmp := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmp)
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}

	var commentIDs []uint64
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Select("id").Where("root_id = ? AND created_at > ?", rootID, cursor).
		Order("created_at DESC").Limit(int(limit)).Pluck("id", &commentIDs).Error
	if err != nil {
		global.Log.Errorf("get child_comment_ids[svc:%v root_id:%v cursor:%v limit:%v] in db failed: %v", svc, rootID, cursor, limit, err)
		return nil, errcode.ERRFindData.WrapError(err)
	}
	return commentIDs, nil
}

func (r *ReadRepo) GetChildCommentCnt(ctx context.Context, svc string, commentID ...uint64) (map[uint64]int, error) {
	if len(commentID) == 0 {
		return nil, nil
	}
	tmp := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmp)
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}

	type Result struct {
		RootID     uint64
		ChildCount int
	}
	var results []Result
	// 查询每个给定 ID 的子评论个数
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).
		Select("root_id, COUNT(*) as child_count").
		Where("root_id IN ?", commentID).
		Group("root_id").
		Scan(&results).Error
	if err != nil {
		global.Log.Errorf("get child_comment_cnt[svc:%v comment_ids:%v] in db failed: %v", svc, commentID, err)
		return nil, errcode.ERRFindData.WrapError(err)
	}

	var mp = make(map[uint64]int, len(results))
	for _, res := range results {
		mp[res.RootID] = res.ChildCount
	}
	return mp, nil
}

func (r *ReadRepo) GetUserIDByCommentID(ctx context.Context, svc string, commentID uint64) (string, error) {
	tmp := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmp)
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return "", errcode.ERRNoTable
	}
	var userID string
	res := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("id = ?", commentID).Pluck("user_id", &userID)
	if res.Error != nil || res.RowsAffected == 0 {
		err := res.Error
		if res.RowsAffected == 0 {
			err = errors.New("get nothing")
		}

		global.Log.Errorf("get user_id from db by comment_id[%v] in svc[%v] failed: %v", commentID, svc, err)

		return "", errcode.ERRFindData.WrapError(res.Error)
	}
	return userID, nil
}

func (r *ReadRepo) CheckCommentExist(ctx context.Context, svc string, commentID ...uint64) map[uint64]bool {
	if len(commentID) == 0 {
		return nil
	}
	tmp := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmp)
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil
	}
	var existingIDs []uint64
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Select("id").Where("id IN ?", commentID).Pluck("id", &existingIDs).Error
	if err != nil {
		global.Log.Errorf("someone is trying to check comment_id[svc:%v comment_id:%v] existence in db failed: %v", svc, commentID, err)
		return nil
	}
	//存储结果，可以只存存在的id
	//不存在的可以不需存，因为查询的时候会查到零值
	result := make(map[uint64]bool, len(commentID))
	for _, id := range existingIDs {
		result[id] = true
	}
	return result
}

func (r *ReadRepo) GetPostList(ctx context.Context, svc string, cursor time.Time, limit uint) ([]uint64, error) {
	tmp := PostInfoPool.Get().(*table.PostInfo)
	defer PostInfoPool.Put(tmp)

	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}

	var postlist []uint64

	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).
		Select("id").Where("created_at < ?", cursor).
		Order("created_at DESC").Limit(int(limit)).
		Pluck("id", &postlist).Error
	if err != nil {
		global.Log.Errorf("get post_list[svc:%v cursor:%v limit:%v] in db failed: %v", svc, cursor, limit, err)
		return nil, errcode.ERRFindData.WrapError(err)
	}

	return postlist, nil
}
