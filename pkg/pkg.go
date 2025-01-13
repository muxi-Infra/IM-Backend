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
