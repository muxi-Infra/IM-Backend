package pg

import (
	"IM-Backend/dao"
	"IM-Backend/errcode"
	"IM-Backend/model/table"
	"context"
	"gorm.io/gorm"
)

type ReadRepo struct {
	db *gorm.DB
	tt dao.TableTooler
}

func (r *ReadRepo) GetPostInfosByExtra(ctx context.Context, svc string, key string, value interface{}) ([]table.PostInfo, error) {
	tmp := &table.PostInfo{}
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}
	tn := tmp.TableName(svc)

	var postInfos []table.PostInfo

	// 使用原生 SQL 查询，根据 extra 字段的 key 查找
	err := r.db.Table(tn).
		Where("extra ->> ? = ?", key, value). // 使用 jsonb 操作符 ->> 提取 key 对应的值并与 value 比较
		Find(&postInfos).Error

	if err != nil {
		return nil, errcode.ERRFindData.WrapError(err)
	}

	return postInfos, nil
}

func (r *ReadRepo) GetCommentsByExtra(ctx context.Context, svc string, key string, value interface{}) ([]table.PostCommentInfo, error) {
	tmp := &table.PostCommentInfo{}
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}
	tn := tmp.TableName(svc)
	var commentInfos []table.PostCommentInfo
	// 使用原生 SQL 查询，根据 extra 字段的 key 查找
	err := r.db.Table(tn).
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

	tmp := &table.PostInfo{}
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}
	tn := tmp.TableName(svc)
	var postInfos []table.PostInfo
	err := r.db.Table(tn).Where("id IN ?", ids).Find(&postInfos).Error
	if err != nil {
		return nil, errcode.ERRFindData.WrapError(err)
	}
	return postInfos, nil
}

func (r *ReadRepo) GetPostLike(ctx context.Context, svc string, id uint64) (int64, error) {
	tmp := &table.PostLikeInfo{}
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return 0, errcode.ERRNoTable
	}
	tn := tmp.TableName(svc)
	var cnt int64
	err := r.db.Table(tn).Where("post_id = ?", id).Count(&cnt).Error
	if err != nil {
		return 0, errcode.ERRCount.WrapError(err)
	}
	return cnt, nil
}

func (r *ReadRepo) GetPostCommentIds(ctx context.Context, svc string, id uint64) ([]uint64, error) {
	tmp := &table.PostCommentInfo{}
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}
	tn := tmp.TableName(svc)
	var commentIds []uint64
	err := r.db.Table(tn).Where("post_id = ?", id).Pluck("id", &commentIds).Error
	if err != nil {
		return nil, errcode.ERRFindData.WrapError(err)
	}
	return commentIds, nil
}

func (r *ReadRepo) GetCommentInfosByID(ctx context.Context, svc string, ids ...uint64) ([]table.PostCommentInfo, error) {
	if len(ids) == 0 {
		return nil, errcode.ERRFindQueryIsEmpty
	}

	tmp := &table.PostCommentInfo{}
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}
	tn := tmp.TableName(svc)
	var commentInfos []table.PostCommentInfo
	err := r.db.Table(tn).Where("id IN ?", ids).Find(&commentInfos).Error
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
	tmp := &table.CommentLikeInfo{}
	if !r.tt.CheckTableExist(r.db, tmp, svc) {
		return nil, errcode.ERRNoTable
	}
	tn := tmp.TableName(svc)

	// 执行查询，按 comment_id 分组，并统计每个 comment_id 的点赞数
	rows, err := r.db.Table(tn).
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
