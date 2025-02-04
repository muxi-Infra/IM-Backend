package service

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"strconv"
	"time"
)

type AuthSvc struct {
	svcHandler SvcHandler
}

func NewAuthSvc(handler SvcHandler) *AuthSvc {
	return &AuthSvc{svcHandler: handler}
}

func (a *AuthSvc) decryptAES(ciphertext string, key []byte) ([]byte, error) {
	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	mode := cipher.NewCFBDecrypter(block, iv)
	mode.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return ciphertextBytes, nil
}

func (a *AuthSvc) Verify(svc string, appKey string) bool {
	secret := a.svcHandler.GetSecretByName(svc)
	decrypted, err := a.decryptAES(appKey, []byte(secret))
	if err != nil {
		return false
	}
	// 将解密后的字节切片转换为字符串
	decryptedStr := string(decrypted)
	// 将字符串解析为整数时间戳
	decryptedTimestamp, err := strconv.ParseInt(decryptedStr, 10, 64)
	if err != nil {
		return false
	}
	// 获取当前时间戳
	currentTimestamp := time.Now().Unix()
	// 计算时间差
	diff := currentTimestamp - decryptedTimestamp
	if diff >= 0 && diff <= 10 {
		return true
	}
	return false
}
