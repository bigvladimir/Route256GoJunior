package utils

import (
	"errors"
	"path/filepath"
	"runtime"
)

const utilsPackageDirDepth = 4

// GetRootPath returns path of root dir
func GetRootPath() (string, error) {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		return "", errors.New("Cant get config path")
	}
	for i := 0; i < utilsPackageDirDepth; i++ {
		path = filepath.Dir(path)
	}
	return path, nil
}
