package controller

import (
	"IM-Backend/controller/req"
	"IM-Backend/controller/resp"
	"IM-Backend/errcode"
	"IM-Backend/model"
	"IM-Backend/model/table"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

type CommentController struct {
	commentSvc CommentService
	ideer      CommentIDGenerator
}

func NewCommentController(commentSvc CommentService, ideer CommentIDGenerator) *CommentController {
	return &CommentController{
		commentSvc: commentSvc,
		ideer:      ideer,
	}
}

func (cc *CommentController) Publish(c *gin.Context) {
	var (
		query    req.PublishCommentQuery
		formdata req.PublishCommentFormData
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	if err := c.ShouldBind(&formdata); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	var extra map[string]interface{}

	if formdata.Extra != nil {
		if err := json.Unmarshal([]byte(*formdata.Extra), &extra); err != nil {
			resp.SendResp(c, resp.ParamBindErrResp)
			return
		}
	}

	createdTime := time.Now()
	id, err := cc.ideer.GenerateCommentID(c, query.Svc)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	var comment = table.PostCommentInfo{
		ID:           id,
		UserID:       query.UserID,
		RootID:       0,
		FatherID:     0,
		TargetUserID: nil,
		PostID:       query.PostID,
		Content:      formdata.Content,
		Extra:        extra,
		CreatedAt:    createdTime,
		UpdatedAt:    createdTime,
	}

	err = cc.commentSvc.Publish(c, query.Svc, comment)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"comment_id": id,
	}))
}

func (cc *CommentController) Reply(c *gin.Context) {
	var (
		query    req.ReplyCommentQuery
		formdata req.ReplyCommentFormData
	)

	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	if query.RootCommentID == 0 || query.FatherCommentID == 0 {
		resp.SendResp(c, resp.ParamErrResp)
		return
	}
	if err := c.ShouldBind(&formdata); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	var extra map[string]interface{}

	if formdata.Extra != nil {
		if err := json.Unmarshal([]byte(*formdata.Extra), &extra); err != nil {
			resp.SendResp(c, resp.ParamBindErrResp)
			return
		}
	}

	userID, err := cc.commentSvc.GetCommentUserIDByID(c, query.Svc, query.FatherCommentID)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}

	id, err := cc.ideer.GenerateCommentID(c, query.Svc)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	createdTime := time.Now()
	var comment = table.PostCommentInfo{
		ID:           id,
		UserID:       query.UserID,
		RootID:       query.RootCommentID,
		FatherID:     query.FatherCommentID,
		TargetUserID: &userID,
		PostID:       query.PostID,
		Content:      formdata.Content,
		Extra:        extra,
		CreatedAt:    createdTime,
		UpdatedAt:    createdTime,
	}

	err = cc.commentSvc.Publish(c, query.Svc, comment)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"comment_id": id,
	}))
}
func (cc *CommentController) Update(c *gin.Context) {
	var (
		query    req.UpdateCommentQuery
		formdata req.UpdateCommentFormData
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	if err := c.ShouldBind(&formdata); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	//如果两个都为空,直接返回
	if formdata.Content == nil && formdata.Extra == nil {
		resp.SendResp(c, resp.NewErrResp(errcode.ERRUpdateQueryEmpty))
		return
	}

	var updates = make(map[string]interface{}, 2)
	if formdata.Content != nil {
		updates["content"] = *formdata.Content
	}
	if formdata.Extra != nil {
		var extra map[string]interface{}
		if err := json.Unmarshal([]byte(*formdata.Extra), &extra); err != nil {
			resp.SendResp(c, resp.ParamBindErrResp)
			return
		}
		updates["extra"] = extra
	}

	err := cc.commentSvc.Update(c, query.Svc, query.UserID, query.CommentID, updates)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	resp.SendResp(c, resp.SuccessResp)
}
func (cc *CommentController) Delete(c *gin.Context) {
	var query req.DeleteCommentQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	err := cc.commentSvc.Delete(c, query.Svc, query.UserID, query.CommentID)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	resp.SendResp(c, resp.SuccessResp)
}

func (cc *CommentController) GetInfo(c *gin.Context) {
	var (
		query req.GetCommentInfoQuery
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	cursor, err := time.Parse("2006-01-02T15:04:05", query.Cursor)
	if err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	comments, err := cc.commentSvc.FindComment(c, query.Svc, query.RootCommentID, cursor, query.Limit)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	res, err := cc.getCommentResp(c, query.Svc, comments)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"comments": res,
	}))
}
func (cc *CommentController) GetLike(c *gin.Context) {
	var (
		query    req.GetCommentLikeQuery
		formdata req.GetCommentLikeJSON
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	if err := c.ShouldBindBodyWithJSON(&formdata); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	if len(formdata.CommentIDs) == 0 {
		resp.SendResp(c, resp.SuccessResp)
		return
	}
	likes, err := cc.commentSvc.GetLike(c, query.Svc, formdata.CommentIDs...)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	var res = make([]int64, 0, len(formdata.CommentIDs))
	for _, id := range formdata.CommentIDs {
		like, ok := likes[id]
		if ok {
			res = append(res, like)
		} else {
			res = append(res, 0)
		}
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"likes": res,
	}))
}
func (cc *CommentController) Like(c *gin.Context) {
	var (
		query req.LikeCommentQuery
		js    req.LikeCommentJson
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	if err := c.ShouldBindBodyWithJSON(&js); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	if js.Like {
		err := cc.commentSvc.Like(c, query.Svc, query.PostID, query.CommentID, query.UserID)
		if err != nil {
			resp.SendResp(c, resp.NewErrResp(err))
			return
		}
	} else {
		err := cc.commentSvc.CancelLike(c, query.Svc, query.CommentID, query.UserID)
		if err != nil {
			resp.SendResp(c, resp.NewErrResp(err))
			return
		}
	}

	resp.SendResp(c, resp.SuccessResp)
}

func (cc *CommentController) getCommentResp(ctx context.Context, svc string, comments []model.PostComment) ([]resp.Comment, error) {
	res := make([]resp.Comment, 0, len(comments))
	ids := make([]uint64, 0, len(comments))
	for _, c := range comments {
		ids = append(ids, c.ID)
	}
	childNums, err := cc.commentSvc.GetChildCommentCnt(ctx, svc, ids...)
	if err != nil {
		return nil, err
	}
	for _, comment := range comments {
		num := childNums[comment.ID]
		if comment.RootID == 0 && comment.FatherID == 0 {
			res = append(res, resp.NewCommentResp(comment, &num))
		} else {
			res = append(res, resp.NewCommentResp(comment, nil))
		}
	}
	return res, nil
}
