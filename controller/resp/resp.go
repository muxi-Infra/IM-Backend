package resp

import (
	"IM-Backend/errcode"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	TimeoutErrResp   = NewResp(http.StatusGatewayTimeout, "Request timeout", nil)
	ParamBindErrResp = NewResp(http.StatusBadRequest, "bind param failed", nil)
	ParamErrResp     = NewResp(http.StatusBadRequest, "something of params is wrong", nil)
	AuthErrResp      = NewResp(http.StatusBadRequest, "bind authentication", nil)
	SuccessResp      = NewResp(http.StatusOK, "success", nil)
)

type Resp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func WithData(r Resp, data interface{}) Resp {
	r.Data = data
	return r
}

func NewResp(code int, msg string, data interface{}) Resp {
	return Resp{
		Code: code,
		Msg:  msg,
		Data: data,
	}
}

func NewErrResp(err error) Resp {
	if err == nil {
		return SuccessResp
	}
	var errc *errcode.Err
	ok := errors.As(err, &errc)
	if !ok {
		return NewResp(http.StatusInternalServerError, "服务端错误", nil)
	}
	return NewResp(errc.Code(), errcode.ClientMsgMapping[errc.Code()], nil)
}

func SendResp(ctx *gin.Context, r Resp) {
	var code int
	if r.Code <= 600 {
		code = r.Code
	} else {
		code = http.StatusInternalServerError
	}
	ctx.JSON(code, r)
}
