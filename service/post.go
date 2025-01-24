package service

import (
	"IM-Backend/configs"
	"IM-Backend/errcode"
	"IM-Backend/model"
	"IM-Backend/model/table"
	"IM-Backend/service/identity"
	"context"
	"errors"
	"time"
)

type PostSvc struct {
	dw          GormWriter
	dr          GormPostReader
	cw          WriteCache
	cr          ReadCache
	postExpire  time.Duration
	pendingPost chan<- identity.PostIdentity
}

func (p *PostSvc) Callback(conf *configs.AppConf) {
	p.postExpire = time.Duration(conf.Cache.PostExpire) * time.Second
}

func (p *PostSvc) Create(ctx context.Context, svc string, postInfo table.PostInfo) error {
	err := p.dw.Create(ctx, svc, &postInfo)
	if err != nil {
		return err
	}
	return nil

}

func (p *PostSvc) Update(ctx context.Context, svc, userID string, postID uint64, updates map[string]interface{}) error {
	oldPostInfo, err := p.getPostInfo(ctx, svc, postID)
	if err != nil && !errors.Is(err, errcode.ERRCacheMiss) {
		return err
	}
	if oldPostInfo.Author != userID {
		return errcode.ERRNoRightRecord
	}
	newPostInfo := oldPostInfo
	tmpCnt := 0
	content, ok := updates["content"]
	if ok {
		contentStr, okk := content.(string)
		if okk {
			tmpCnt++
			newPostInfo.Content = contentStr
		}
	}

	extra, ok := updates["extra"]
	if ok {
		extraMap, okk := extra.(map[string]interface{})
		if okk {
			tmpCnt++
			newPostInfo.Extra = extraMap
		}
	}
	if tmpCnt == 0 {
		return errcode.ERRUpdateQueryEmpty
	}
	newPostInfo.UpdatedAt = time.Now()
	if err := p.cw.DelKV(ctx, &oldPostInfo); err != nil {
		return err
	}
	tmp := newPostInfo.ToTable()
	if err := p.dw.Update(ctx, svc, &tmp); err != nil {
		return err
	}
	go func() {
		time.AfterFunc(time.Second, func() {
			_ = p.cw.DelKV(context.Background(), &oldPostInfo)
		})
	}()

	return nil
}

func (p *PostSvc) GetInfo(ctx context.Context, svc string, postID uint64) (model.PostInfo, error) {
	res, err := p.getPostInfo(ctx, svc, postID)
	if err == nil { //缓存命中
		return res, nil
	}
	if errors.Is(err, errcode.ERRCacheMiss) { //缓存未命中
		_ = p.cw.SetKV(ctx, p.postExpire, &res)
		return res, nil
	}
	return model.PostInfo{}, err
}

func (p *PostSvc) GetLike(ctx context.Context, svc string, postID uint64) (int64, error) {
	res, err := p.getLike(ctx, svc, postID)
	if err == nil {
		return res, nil
	}
	if errors.Is(err, errcode.ERRCacheMiss) {
		_ = p.cw.SetKV(ctx, 5*time.Minute, &model.PostLike{
			PostID: postID,
			Svc:    svc,
			Like:   res,
		})
		return res, nil
	}
	return 0, err
}

func (p *PostSvc) Delete(ctx context.Context, svc string, userID string, postID uint64) error {
	//检查是否合法
	res, err := p.getPostInfo(ctx, svc, postID)
	if err != nil && !errors.Is(err, errcode.ERRCacheMiss) {
		return err
	}
	if res.Author != userID {
		return errcode.ERRNoRightRecord
	}

	//删除缓存
	go func() {
		_ = p.cw.DelKV(context.Background(), &model.PostLike{
			PostID: postID,
			Svc:    svc,
		})
	}()

	if err := p.cw.DelKV(ctx, &model.PostInfo{
		ID:  postID,
		Svc: svc,
	}); err != nil {
		return err
	}
	//删除主要信息
	if err := p.dw.Delete(ctx, svc, &table.PostInfo{ID: postID}); err != nil {
		return err
	}

	//去通知清理协程异步地清理相关痕迹(如评论，点赞)
	p.pendingPost <- identity.PostIdentity{
		Svc:    svc,
		PostID: postID,
	}

	return nil
}

func (p *PostSvc) Like(ctx context.Context, svc string, postID uint64, userID string) error {
	//删除点赞缓存(根据点赞数选择相应策略)
	if err := p.delLike(ctx, svc, postID); err != nil {
		return err
	}

	//保存点赞记录
	if err := p.dw.Create(ctx, svc, &table.PostLikeInfo{
		PostID:    postID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}); err != nil {
		return err
	}

	go func() {
		time.AfterFunc(time.Second, func() {
			_ = p.delLike(context.Background(), svc, postID)
		})
	}()

	return nil
}

func (p *PostSvc) getPostInfo(ctx context.Context, svc string, postID uint64) (model.PostInfo, error) {
	var res = model.PostInfo{ID: postID, Svc: svc}
	err := p.cr.GetKV(ctx, &res)
	if err == nil {
		return res, nil
	}
	postInfos, err := p.dr.GetPostInfos(ctx, svc, postID)
	if err != nil {
		return model.PostInfo{}, err
	}
	if len(postInfos) != 1 {
		return model.PostInfo{}, errcode.ERRNoRightRecord
	}
	return model.NewPostInfo(postInfos[0], svc), errcode.ERRCacheMiss
}

func (p *PostSvc) getLike(ctx context.Context, svc string, postID uint64) (int64, error) {
	//先检查是否存在该post
	if !p.dr.CheckExist(ctx, svc, postID) {
		return 0, nil
	}

	var res = model.PostLike{
		PostID: postID,
		Svc:    svc,
	}
	err := p.cr.GetKV(ctx, &res)
	if err == nil {
		return res.Like, nil
	}
	cnt, err := p.dr.GetPostLike(ctx, svc, postID)
	if err != nil {
		return 0, err
	}
	return cnt, errcode.ERRCacheMiss
}

func (p *PostSvc) delLike(ctx context.Context, svc string, postID uint64) error {
	var tmp = &model.PostLike{
		PostID: postID,
		Svc:    svc,
	}
	err := p.cr.GetKV(ctx, tmp)
	//当查询点赞数小于10000或者查询失败则删除缓存的点赞数
	//其他情况则不删，等待过期时间结束
	if err != nil || tmp.Like < 10000 {
		if err := p.cw.DelKV(ctx, tmp); err != nil {
			return err
		}
	}
	return nil
}
