package service

import (
	"IM-Backend/cache"
	"IM-Backend/configs"
	"IM-Backend/errcode"
	"IM-Backend/model"
	"IM-Backend/model/table"
	"IM-Backend/pkg"
	"IM-Backend/service/identity"
	"context"
	"errors"
	"log"
	"time"
)

type CommentSvc struct {
	dw            GormWriter
	dr            GormCommentReader
	cw            WriteCache
	cr            ReadCache
	commentExpire time.Duration

	pendingComments chan<- identity.CommentIdentity
}

func NewCommentService(dw GormWriter, dr GormCommentReader, cw WriteCache, cr ReadCache, pendingComments chan<- identity.CommentIdentity, cf configs.AppConf) *CommentSvc {
	return &CommentSvc{
		dw:              dw,
		dr:              dr,
		cw:              cw,
		cr:              cr,
		commentExpire:   time.Duration(cf.Cache.CommentExpire) * time.Second,
		pendingComments: pendingComments,
	}
}

func (c *CommentSvc) Callback(conf configs.AppConf) {
	c.commentExpire = time.Duration(conf.Cache.CommentExpire) * time.Second
	log.Printf("comment expire has been changed to %v", c.commentExpire)
}

func (c *CommentSvc) GetCommentUserIDByID(ctx context.Context, svc string, commentID uint64) (string, error) {
	//可以尝试看看缓存中有没有comment的信息
	//如果有可以直接返回
	var comment = model.PostComment{
		ID:  commentID,
		Svc: svc,
	}
	err := c.cr.GetKV(ctx, &comment)
	if err == nil {
		return comment.UserID, nil
	}
	//如果没有,则直接查询数据库
	//也无需考虑缓存，因为只获取userID
	return c.dr.GetUserIDByCommentID(ctx, svc, commentID)
}

func (c *CommentSvc) Like(ctx context.Context, svc string, postID uint64, commentID uint64, userID string) error {
	//点赞需要删除之前的缓存的点赞信息
	if err := c.delLike(ctx, svc, commentID); err != nil {
		return err
	}
	err := c.dw.Create(ctx, svc, &table.CommentLikeInfo{
		CommentID: commentID,
		UserID:    userID,
		PostID:    postID,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return err
	}
	//延迟删除
	go func() {
		time.AfterFunc(time.Second, func() {
			_ = c.delLike(context.Background(), svc, commentID)
		})
	}()
	return nil
}

// Publish 发表评论
func (c *CommentSvc) Publish(ctx context.Context, svc string, comment table.PostCommentInfo) error {
	err := c.dw.Create(ctx, svc, &comment)
	if err != nil {
		return err
	}
	return nil
}

// Update 更新评论信息
func (c *CommentSvc) Update(ctx context.Context, svc, userID string, commentID uint64, updates map[string]interface{}) error {
	//获取老的comment
	var oldComment model.PostComment
	found, left, err := c.getCommentsByID(ctx, svc, commentID)
	if err != nil && !errors.Is(err, errcode.ERRCacheMiss) {
		return err
	}
	if err == nil {
		//说明在缓存中找到了
		oldComment = found[0]
	} else {
		//缓存未命中
		oldComment = left[0]
	}

	//验证userid是否配对
	if oldComment.UserID != userID {
		return errcode.ERRNoRightRecord
	}
	//对评论信息进行更新
	newComment, err := c.updateCommentFromMap(oldComment, updates)
	if err != nil {
		return err
	}

	//删除原来评论的缓存
	if err := c.cw.DelKV(ctx, &oldComment); err != nil {
		return err
	}
	tmp := newComment.ToTable()
	if err = c.dw.Update(ctx, svc, &tmp); err != nil {
		return err
	}

	//延迟删除缓存
	go func() {
		time.AfterFunc(time.Second, func() {
			_ = c.cw.DelKV(context.Background(), &oldComment)
		})
	}()
	return nil
}
func (c *CommentSvc) FindComment(ctx context.Context, svc string, rootID uint64, cursor time.Time, limit uint) ([]model.PostComment, error) {
	//首先,先从数据库中获取root_id == rootID的评论(且创建时间大于cursor同时限制个数为limit)有哪些
	childCommentIDs, err := c.dr.GetChildCommentIDAfterCursor(ctx, svc, rootID, cursor, limit)
	if err != nil {
		//获取失败，则返回
		return nil, err
	}
	//没有符合条件的，直接返回，但不返回错误
	if len(childCommentIDs) == 0 {
		return nil, nil
	}

	found, left, err := c.getCommentsByID(ctx, svc, childCommentIDs...)
	if err == nil { //缓存全部命中
		return found, nil
	}
	if !errors.Is(err, errcode.ERRCacheMiss) { //查询失败
		return nil, err
	}
	go func() { //缓存未全部命中
		//进行缓存
		_ = c.setCommentCache(context.Background(), left...)
	}()
	//返回已经命中的部分和未命中的部分
	return append(found, left...), nil
}

func (c *CommentSvc) GetLike(ctx context.Context, svc string, commentID ...uint64) (map[uint64]int64, error) {
	found, left, err := c.getLike(ctx, svc, commentID...)
	if err == nil {
		return found, nil
	}
	if !errors.Is(err, errcode.ERRCacheMiss) {
		return nil, err
	}
	go func() {
		//进行缓存
		_ = c.setLikeCache(context.Background(), svc, left)
	}()
	return pkg.MergeMaps(found, left), nil
}

func (c *CommentSvc) GetChildCommentCnt(ctx context.Context, svc string, commentID ...uint64) (map[uint64]int, error) {
	if len(commentID) == 0 {
		return nil, nil
	}
	return c.dr.GetChildCommentCnt(ctx, svc, commentID...)
}

func (c *CommentSvc) Delete(ctx context.Context, svc string, userID string, commentID uint64) error {
	//检查要删除的是否合法
	found, left, err := c.getCommentsByID(ctx, svc, commentID)
	if err != nil && !errors.Is(err, errcode.ERRCacheMiss) {
		return err
	}

	if err == nil {
		if found[0].UserID != userID {
			return errcode.ERRNoRightRecord
		}
	} else {
		if left[0].UserID != userID {
			return errcode.ERRNoRightRecord
		}
	}

	//删除缓存
	go func() {
		_ = c.cw.DelKV(context.Background(), &model.PostCommentLike{CommentID: commentID, Svc: svc})
	}()

	if err := c.cw.DelKV(ctx, &model.PostComment{ID: commentID, Svc: svc}); err != nil {
		return err
	}

	//删除基本信息
	if err := c.dw.Delete(ctx, svc, &table.PostCommentInfo{ID: commentID}); err != nil {
		return err
	}
	c.pendingComments <- identity.CommentIdentity{
		Svc:       svc,
		CommentID: commentID,
	}
	return nil

}

func (c *CommentSvc) getCommentsByIDFromCache(ctx context.Context, svc string, commentID ...uint64) ([]model.PostComment /*查询到的*/, []uint64 /*未查询到的ID*/) {
	if len(commentID) == 0 {
		return nil, nil
	}
	var (
		kvs    = make([]cache.KV, len(commentID))
		found  = make([]model.PostComment, 0)
		leftID = make([]uint64, 0)
	)
	for i, id := range commentID {
		kvs[i] = &model.PostComment{ID: id, Svc: svc}
	}
	res := c.cr.MGetKV(ctx, kvs...)
	for i, ok := range res {
		if ok {
			found = append(found, *(kvs[i].(*model.PostComment)))
		} else {
			leftID = append(leftID, commentID[i])
		}
	}
	return found, leftID
}

// 如果查询成功
// 返回在缓存中查到的和未查到的
func (c *CommentSvc) getCommentsByID(ctx context.Context, svc string, commentID ...uint64) ([]model.PostComment, []model.PostComment, error) {
	//有了id可以查询缓存
	found, leftID := c.getCommentsByIDFromCache(ctx, svc, commentID...)
	//如果没有未被查询的,则直接返回
	if len(leftID) == 0 {
		return found, nil, nil
	}
	//接下来从数据库中查询剩余的
	leftComments, err := c.dr.GetCommentInfosByID(ctx, svc, leftID...)
	if err != nil {
		//数据库查询失败
		return nil, nil, err
	}

	//left为缓存未命中的部分
	left := make([]model.PostComment, 0)

	for _, comment := range leftComments {
		left = append(left, model.NewPostComment(comment, svc))
	}
	return found, left, errcode.ERRCacheMiss
}

// 缓存
func (c *CommentSvc) setCommentCache(ctx context.Context, comments ...model.PostComment) error {
	if len(comments) == 0 {
		return nil
	}
	var kvs = make([]cache.KV, len(comments))
	for i, comment := range comments {
		kvs[i] = &comment
	}
	return c.cw.SetKV(ctx, c.commentExpire, kvs...)
}

func (c *CommentSvc) getLikeFromCache(ctx context.Context, svc string, commentID ...uint64) (map[uint64]int64, []uint64) {
	if len(commentID) == 0 {
		return nil, nil
	}
	var (
		kvs    = make([]cache.KV, len(commentID))
		found  = make(map[uint64]int64)
		leftID = make([]uint64, 0)
	)
	for i, id := range commentID {
		kvs[i] = &model.PostCommentLike{
			Svc:       svc,
			CommentID: id,
		}
	}
	res := c.cr.MGetKV(ctx, kvs...)
	for i, ok := range res {
		if ok {
			found[commentID[i]] = (kvs[i].(*model.PostCommentLike)).Like
		} else {
			leftID = append(leftID, commentID[i])
		}
	}
	return found, leftID
}
func (c *CommentSvc) getLike(ctx context.Context, svc string, commentID ...uint64) (map[uint64]int64, map[uint64]int64, error) {
	if len(commentID) == 0 {
		return nil, nil, nil
	}
	//如果没有某个comment则不查询其
	var unexist = make(map[uint64]int64) //unexist是用来存储不存在的
	var rightCommentID []uint64          //rightCommentID是用来存储存在的commentID
	exist := c.dr.CheckCommentExist(ctx, svc, commentID...)
	for id, ok := range exist {
		if ok {
			//这里是存在的comment
			rightCommentID = append(rightCommentID, id)
		} else {
			//不存在的直接置为0，无需查询
			unexist[id] = 0
		}
	}

	//如果全是不存在的，直接返回unexist即可
	if len(rightCommentID) == 0 {
		return unexist, nil, nil
	}

	found, leftID := c.getLikeFromCache(ctx, svc, rightCommentID...)
	//全部在缓存中找到
	if len(leftID) == 0 {
		return pkg.MergeMaps(found, unexist), nil, nil
	}
	res, err := c.dr.GetCommentLike(ctx, svc, leftID...)

	if err != nil {
		return nil, nil, err
	}

	return pkg.MergeMaps(found, unexist), res, errcode.ERRCacheMiss
}

func (c *CommentSvc) setLikeCache(ctx context.Context, svc string, mp map[uint64]int64) error {
	if len(mp) == 0 {
		return nil
	}
	var kvs = make([]cache.KV, 0, len(mp))
	for id, like := range mp {
		kvs = append(kvs, &model.PostCommentLike{
			Svc:       svc,
			CommentID: id,
			Like:      like,
		})
	}
	return c.cw.SetKV(ctx, 5*time.Minute, kvs...)
}

func (c *CommentSvc) delLike(ctx context.Context, svc string, commentID uint64) error {
	var tmp = &model.PostCommentLike{
		Svc:       svc,
		CommentID: commentID,
	}
	err := c.cr.GetKV(ctx, tmp)
	//当查询点赞数小于10000或者查询失败则删除缓存的点赞数
	//其他情况则不删，等待过期时间结束
	if err != nil || tmp.Like < 10000 {
		if err := c.cw.DelKV(ctx, tmp); err != nil {
			return err
		}
	}
	return nil
}

func (c *CommentSvc) updateCommentFromMap(oldComment model.PostComment, updates map[string]interface{}) (model.PostComment, error) {
	//对评论信息进行更新
	newComment := oldComment
	content, ok := updates["content"]
	tmpCnt := 0 //查看是否更新了
	if ok {
		contentStr, okk := content.(string)
		if okk {
			tmpCnt++
			newComment.Content = contentStr
		}
	}
	extra, ok := updates["extra"]
	if ok {
		extraMap, okk := extra.(map[string]interface{})
		if okk {
			tmpCnt++
			newComment.Extra = extraMap
		}
	}
	//未更新任何内容，直接返回
	if tmpCnt == 0 {
		return model.PostComment{}, errcode.ERRUpdateQueryEmpty
	}
	newComment.UpdatedAt = time.Now()
	return newComment, nil
}
