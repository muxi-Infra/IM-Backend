package route

import (
	"IM-Backend/controller"
	"IM-Backend/middleware"
	"IM-Backend/service"
	"context"

	"github.com/gin-gonic/gin"
)

type App struct {
	r        *gin.Engine
	dectSvc  *service.DetectSvc
	cleanSvc *service.CleanSvc
}

func NewApp(pc *controller.PostController, cc *controller.CommentController, authSvc *service.AuthSvc,
	dectSvc *service.DetectSvc, cleanSvc *service.CleanSvc) *App {
	r := gin.Default()
	//加载中间件
	loadMiddleware(r, middleware.LockMiddleware(),
		middleware.TimeoutMiddleware(),
		middleware.AuthMiddleware(authSvc))
	loadRoute(r, pc, cc)
	return &App{
		r:        r,
		dectSvc:  dectSvc,
		cleanSvc: cleanSvc,
	}
}

func loadRoute(r *gin.Engine, pc *controller.PostController, cc *controller.CommentController) {

	postGroup := r.Group("/api/v1/posts")
	{
		postGroup.POST("/publish", pc.Publish)
		postGroup.GET("/getinfo", pc.GetInfo)
		postGroup.PUT("/update", pc.Update)
		postGroup.DELETE("/delete", pc.Delete)
		postGroup.PUT("/like", pc.Like)
		postGroup.GET("/getlike", pc.GetLike)
		postGroup.GET("/getlist", pc.GetList)
	}
	postCommentGroup := postGroup.Group("/comments")
	{
		postCommentGroup.POST("/publish", cc.Publish)
		postCommentGroup.POST("/reply", cc.Reply)
		postCommentGroup.PUT("/update", cc.Update)
		postCommentGroup.DELETE("/delete", cc.Delete)
		postCommentGroup.GET("/getinfo", cc.GetInfo)
		postCommentGroup.PUT("/like", cc.Like)
		postCommentGroup.GET("/getlike", cc.GetLike)
	}
}

func loadMiddleware(r *gin.Engine, middleware ...gin.HandlerFunc) {
	for _, m := range middleware {
		r.Use(m)
	}
}

func (a *App) Run(ctx context.Context) {
	go a.dectSvc.Run(ctx)
	go a.cleanSvc.Run(ctx)
	if err := a.r.Run(":8181"); err != nil {
		panic(err)
	}
}
