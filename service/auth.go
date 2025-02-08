package service

import (
	"IM-Backend/global"
	"IM-Backend/pkg"
	"strconv"
	"time"
)

type AuthSvc struct {
	svcHandler SvcHandler
}

func NewAuthSvc(handler SvcHandler) *AuthSvc {
	return &AuthSvc{svcHandler: handler}
}

func (a *AuthSvc) Verify(svc string, appKey string) bool {
	secret := a.svcHandler.GetSecretByName(svc)
	decrypted, err := pkg.DecryptAES(appKey, []byte(secret))
	if err != nil {
		return false
	}
	// 将解密后的字节切片转换为字符串
	decryptedStr := string(decrypted)
	// 将字符串解析为整数时间戳
	decryptedTimestamp, err := strconv.ParseInt(decryptedStr, 10, 64)
	if err != nil {
		global.Log.Errorf("parse appKey[%s] failed: %v", appKey, err)
		return false
	}
	// 获取当前时间戳
	currentTimestamp := time.Now().Unix()
	// 计算时间差

	diff := currentTimestamp - decryptedTimestamp

	if diff < 0 {
		global.Log.Warnf("time may be slower than right time,diff: %v", diff)
	}

	if diff >= -10 && diff <= 10 {
		return true
	}
	return false
}
