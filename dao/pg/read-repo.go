package pg

import (
	"IM-Backend/dao"
	"IM-Backend/errcode"
	"IM-Backend/model/table"
	"context"
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
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("id IN ?", ids).Find(&postInfos).Error
	if err != nil {
		return nil, errcode.ERRFindData.WrapError(err)
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
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("id IN ?", ids).Find(&commentInfos).Error
	if err != nil {
		return nil, errcode.ERRFindData.WrapError(err)
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
		return nil, errcode.ERRCount.WrapError(err)
	}

	return result, nil
}

func (r *ReadRepo) GetChildCommentIDAfterCursor(ctx context.Context, svc string, fatherID uint64, cursor time.Time, limit uint) ([]uint64, error) {
	tmp := CommentInfoPool.Get().(*table.PostCommentInfo)
	defer CommentInfoPool.Put(tmp)
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}

	var commentIDs []uint64
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Select("id").Where("father_id = ? AND created_at > ?", fatherID, cursor).
		Order("created_at").Limit(int(limit)).Pluck("id", &commentIDs).Error
	if err != nil {
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
		FatherID   uint64
		ChildCount int
	}
	var results []Result
	// 查询每个给定 ID 的子评论个数
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).
		Select("father_id, COUNT(*) as child_count").
		Where("father_id IN ?", commentID).
		Group("father_id").
		Scan(&results).Error
	if err != nil {
		return nil, errcode.ERRFindData.WrapError(err)
	}

	var mp = make(map[uint64]int, len(results))
	for _, res := range results {
		mp[res.FatherID] = res.ChildCount
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
	err := r.db.WithContext(ctx).Table(tmp.TableName(svc)).Where("id = ?", commentID).Pluck("user_id", &userID).Error
	if err != nil {
		return "", errcode.ERRFindData.WrapError(err)
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
