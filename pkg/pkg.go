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
