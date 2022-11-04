package util

import (
	"errors"
	"path/filepath"
	"runtime"
)

func CodeDirName() (string, error) {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "", errors.New("unable to get code directory name")
	}
	return filepath.Dir(filename), nil
}
