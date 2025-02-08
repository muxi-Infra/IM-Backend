package service

import (
	"IM-Backend/configs"
	"IM-Backend/errcode"
	"IM-Backend/global"
	"IM-Backend/model"
	"IM-Backend/model/table"
	"IM-Backend/service/identity"
	"context"
	"errors"
	"time"
)

type PostSvc struct {
	dw             GormWriter
	dr             GormPostReader
	cw             WriteCache
	cr             ReadCache
	postExpire     time.Duration
	postLikeExpire time.Duration
	pendingPost    chan<- identity.PostIdentity
}

func NewPostSvc(dw GormWriter, dr GormPostReader, cw WriteCache, cr ReadCache,
	pendingPost chan<- identity.PostIdentity, cf configs.AppConf) *PostSvc {
	return &PostSvc{
		dw:             dw,
		dr:             dr,
		cw:             cw,
		cr:             cr,
		pendingPost:    pendingPost,
		postExpire:     time.Duration(cf.Cache.PostExpire) * time.Second,
		postLikeExpire: time.Duration(cf.Cache.PostLikeExpire) * time.Second,
	}
}

func (p *PostSvc) Callback(conf configs.AppConf) {
	p.postExpire = time.Duration(conf.Cache.PostExpire) * time.Second
	p.postLikeExpire = time.Duration(conf.Cache.PostLikeExpire) * time.Second
	global.Log.Infof("post and postlike expire has been changed to [%v,%v]", p.postExpire, p.postLikeExpire)
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
		global.Log.Errorf("update post[id:%v] in svc[%v],but it's author != %v", postID, svc, userID)
		return errcode.ERRNoRightRecord
	}

	newPostInfo, err := p.updatePostInfoFromMap(oldPostInfo, updates)
	if err != nil {
		return err
	}

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
		_ = p.cw.SetKV(ctx, p.postLikeExpire, &model.PostLike{
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
		global.Log.Errorf("delete post[id:%v] in svc[%v],but it's author != %v", postID, svc, userID)
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
	if err := p.dw.Delete(ctx, svc, &table.PostInfo{ID: postID}, nil); err != nil {
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

	//延迟删除缓存
	go func() {
		time.AfterFunc(time.Second, func() {
			_ = p.delLike(context.Background(), svc, postID)
		})
	}()

	return nil
}

func (p *PostSvc) CancelLike(ctx context.Context, svc string, postID uint64, userID string) error {
	//删除点赞缓存(根据点赞数选择相应策略)
	if err := p.delLike(ctx, svc, postID); err != nil {
		return err
	}

	if err := p.dw.Delete(ctx, svc, &table.PostLikeInfo{}, map[string]interface{}{
		"post_id": postID,
		"user_id": userID,
	}); err != nil {
		return err
	}

	//延迟删除缓存
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

	return model.NewPostInfo(postInfos[0], svc), errcode.ERRCacheMiss
}

func (p *PostSvc) getLike(ctx context.Context, svc string, postID uint64) (int64, error) {
	//先检查是否存在该post
	if !p.dr.CheckPostExist(ctx, svc, postID) {
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

func (p *PostSvc) updatePostInfoFromMap(oldPostInfo model.PostInfo, updates map[string]interface{}) (model.PostInfo, error) {
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
		return model.PostInfo{}, errcode.ERRUpdateQueryEmpty
	}
	newPostInfo.UpdatedAt = time.Now()
	return newPostInfo, nil
}
