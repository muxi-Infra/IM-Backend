package middleware

import (
	"IM-Backend/controller/resp"
	"IM-Backend/global"
	"IM-Backend/service"
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
)

// 在处理请求前加锁，在处理请求后解锁
func LockMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 加锁
		global.AppLock.RLock()
		defer global.AppLock.RUnlock()

		// 继续处理请求
		c.Next()
	}
}

// 超时中间件
func TimeoutMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 创建带有超时的上下文
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel() // 确保在函数结束时取消上下文

		// 替换请求的上下文为带超时的版本
		c.Request = c.Request.WithContext(ctx)

		// 创建一个通道用于监控请求是否完成
		done := make(chan struct{})

		go func() {
			defer close(done)
			c.Next() // 继续处理后续中间件和路由处理函数
		}()

		select {
		case <-done:
			// 请求正常完成
			return
		case <-ctx.Done():
			// 超时触发
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				resp.SendResp(c, resp.TimeoutErrResp)
				c.Abort()
				return
			}
		}
	}
}

// 鉴权中间件
func AuthMiddleware(authSvc *service.AuthSvc) gin.HandlerFunc {
	return func(c *gin.Context) {
		type tmp struct {
			Svc    string `form:"svc"`
			AppKey string `form:"appKey"`
		}
		var t tmp
		//绑定参数
		if err := c.ShouldBindQuery(&t); err != nil {
			//直接返回错误响应
			resp.SendResp(c, resp.ParamBindErrResp)
			c.Abort()
			return
		}
		//验证
		ok := authSvc.Verify(t.Svc, t.AppKey)
		if !ok {
			//返回错误响应
			resp.SendResp(c, resp.AuthErrResp)
			c.Abort()
			return
		}
		c.Next()
	}
}
