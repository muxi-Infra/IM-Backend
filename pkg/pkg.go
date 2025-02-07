package pkg

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

func ReadYaml(path string, aim any) error {
	dataBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = ReadYamlContent(dataBytes, aim)
	return nil
}

func ReadYamlContent(content []byte, aim any) error {
	err := yaml.Unmarshal(content, aim)
	if err != nil {
		return err
	}
	return nil
}

// MergeMaps 合并两个 map 的函数
func MergeMaps[K comparable, V any](map1, map2 map[K]V) map[K]V {
	if len(map1) == 0 {
		return map2
	}
	if len(map2) == 0 {
		return map1
	}

	result := make(map[K]V)

	// 将第一个 map 中的所有键值对复制到结果 map 中
	for k, v := range map1 {
		result[k] = v
	}

	// 将第二个 map 中的所有键值对添加到结果 map 中
	for k, v := range map2 {
		result[k] = v
	}

	return result
}

// Unique 通用去重函数，适用于任意可比较类型
func Unique[T comparable](slice []T) []T {
	// 用于记录元素是否已经存在
	seen := make(map[T]struct{})
	// 存储去重后的结果
	var result = make([]T, 0)

	// 遍历输入切片
	for _, item := range slice {
		// 如果元素还没有出现过
		if _, ok := seen[item]; !ok {
			// 将元素添加到结果切片中
			result = append(result, item)
			// 标记该元素已经出现过
			seen[item] = struct{}{}
		}
	}
	return result
}

// DecryptAES 使用 AES 的 CFB 模式解密
func DecryptAES(ciphertext string, key []byte) ([]byte, error) {
	ciphertextBytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// 使用前 AES BlockSize 的长度作为 IV
	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]

	mode := cipher.NewCFBDecrypter(block, iv)
	mode.XORKeyStream(ciphertextBytes, ciphertextBytes)

	return ciphertextBytes, nil
}

// EncryptAES 使用 AES 的 CFB 模式对明文进行加密，并将结果进行 Hex 编码
func EncryptAES(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 生成 IV
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return "", err
	}

	// 使用 CFB 加密模式
	stream := cipher.NewCFBEncrypter(block, iv)
	ciphertext := make([]byte, len(plaintext))
	stream.XORKeyStream(ciphertext, []byte(plaintext))

	// 将 IV 和加密后的数据拼接在一起
	ciphertextWithIV := append(iv, ciphertext...)

	// 返回 Hex 编码的密文
	return hex.EncodeToString(ciphertextWithIV), nil
}

// FormatTimeInShanghai 将 time.Time 转换为 Asia/Shanghai 时区的格式化时间字符串
func FormatTimeInShanghai(t time.Time) string {
	// 获取 Shanghai 时区
	location, _ := time.LoadLocation("Asia/Shanghai")

	// 将时间转换为 Asia/Shanghai 时区
	shanghaiTime := t.In(location)

	// 格式化并返回
	return shanghaiTime.Format("2006-01-02T15:04")
}
