package kafka

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"

	"homework/internal/pkg/utils"
)

func GetBrokers(path string) ([]string, error) {
	rootPath, err := utils.GetRootPath()
	if err != nil {
		return nil, fmt.Errorf("utils.GetRootPath: %w", err)
	}
	path = filepath.Join(rootPath, filepath.Clean(path))
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("io.ReadFile: %w", err)
	}

	data := make(map[string]string)
	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	// Преобразование словаря в слайс строк
	var brokers []string
	for _, value := range data {
		brokers = append(brokers, value)
	}

	return brokers, nil
}
