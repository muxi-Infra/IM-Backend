package pkg

import (
	"os"

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
	var result []T

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
