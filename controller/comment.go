package controller

import (
	"IM-Backend/controller/req"
	"IM-Backend/controller/resp"
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
	if err := json.Unmarshal([]byte(formdata.Extra), &extra); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
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
		FatherID:     0,
		TargetUserID: "none",
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
	if err := c.ShouldBind(&formdata); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	var extra map[string]interface{}
	if err := json.Unmarshal([]byte(formdata.Extra), &extra); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	createdTime := time.Now()

	userID, err := cc.commentSvc.GetCommentUserIDByID(c, query.Svc, query.CommentID)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}

	id, err := cc.ideer.GenerateCommentID(c, query.Svc)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}

	var comment = table.PostCommentInfo{
		ID:           id,
		UserID:       query.UserID,
		FatherID:     query.CommentID,
		TargetUserID: userID,
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
	var extra map[string]interface{}
	if err := json.Unmarshal([]byte(formdata.Extra), &extra); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	err := cc.commentSvc.Update(c, query.Svc, query.UserID, query.CommentID, map[string]interface{}{
		"content": formdata.Content,
		"extra":   extra,
	})
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
	comments, err := cc.commentSvc.FindComment(c, query.Svc, query.FatherCommentID, cursor, query.Limit)
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
		formdata req.GetCommentLikeFormData
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	if err := c.ShouldBind(&formdata); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
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
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	err := cc.commentSvc.Like(c, query.Svc, query.PostID, query.CommentID, query.UserID)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
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
		res = append(res, resp.NewCommentResp(comment, childNums[comment.ID]))
	}
	return res, nil

}
