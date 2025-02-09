package controller

import (
	"IM-Backend/controller/req"
	"IM-Backend/controller/resp"
	"IM-Backend/errcode"
	"IM-Backend/model/table"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"time"
)

type PostController struct {
	postSvc PostService
	ider    PostIDGenerator
}

func NewPostController(postSvc PostService, ider PostIDGenerator) *PostController {
	return &PostController{
		postSvc: postSvc,
		ider:    ider,
	}
}

func (p *PostController) Publish(c *gin.Context) {
	var (
		query    req.PublishPostQuery
		formdata req.PublishPostFormData
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

	var post = table.PostInfo{
		Author:    formdata.Author,
		Title:     formdata.Title,
		Content:   formdata.Content,
		Extra:     extra,
		CreatedAt: createdTime,
		UpdatedAt: createdTime,
	}
	id, err := p.ider.GeneratePostID(c, query.Svc)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	post.ID = id

	err = p.postSvc.Create(c, query.Svc, post)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"post_id": id,
	}))
}

func (p *PostController) GetList(c *gin.Context) {
	var query req.GetPostListQuery

	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	cursor, err := time.Parse("2006-01-02T15:04:05", query.Cursor)
	if err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	ids, err := p.postSvc.GetList(c, query.Svc, cursor, query.Limit)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}

	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"post_list": ids,
	}))
}

func (p *PostController) GetInfo(c *gin.Context) {
	var query req.GetPostInfoQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	postInfo, err := p.postSvc.GetInfo(c, query.Svc, query.PostID)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"post": resp.NewPostResp(postInfo),
	}))
}
func (p *PostController) Delete(c *gin.Context) {
	var query req.DeletePostQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	err := p.postSvc.Delete(c, query.Svc, query.UserID, query.PostID)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	resp.SendResp(c, resp.SuccessResp)
}
func (p *PostController) GetLike(c *gin.Context) {
	var query req.GetPostLikeQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	likeCount, err := p.postSvc.GetLike(c, query.Svc, query.PostID)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	resp.SendResp(c, resp.WithData(resp.SuccessResp, map[string]interface{}{
		"like": likeCount,
	}))
}

func (p *PostController) Like(c *gin.Context) {
	var query req.LikePostQuery
	var js req.LikePostJson
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	if err := c.ShouldBindJSON(&js); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}

	if js.Like {
		err := p.postSvc.Like(c, query.Svc, query.PostID, query.UserID)
		if err != nil {
			resp.SendResp(c, resp.NewErrResp(err))
			return
		}
	} else {
		err := p.postSvc.CancelLike(c, query.Svc, query.PostID, query.UserID)
		if err != nil {
			resp.SendResp(c, resp.NewErrResp(err))
			return
		}
	}

	resp.SendResp(c, resp.SuccessResp)
}

func (p *PostController) Update(c *gin.Context) {
	var (
		query    req.UpdatePostQuery
		formdata req.UpdatePostFormData
	)
	if err := c.ShouldBindQuery(&query); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	if err := c.ShouldBind(&formdata); err != nil {
		resp.SendResp(c, resp.ParamBindErrResp)
		return
	}
	//如果均为空则直接返回
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

	err := p.postSvc.Update(c, query.Svc, query.UserID, query.PostID, updates)
	if err != nil {
		resp.SendResp(c, resp.NewErrResp(err))
		return
	}
	resp.SendResp(c, resp.SuccessResp)
}
