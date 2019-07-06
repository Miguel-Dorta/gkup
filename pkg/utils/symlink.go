package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func IsSymLink(mode os.FileMode) bool {
	return mode&os.ModeSymlink != 0
}

func ResolveSymlink(path string) (os.FileInfo, error) {
	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return nil, fmt.Errorf("cannot eval symlink in \"%s\": %s", path, err.Error())
	}

	stat, err := os.Stat(realPath)
	if err != nil {
		return nil, fmt.Errorf("cannot get file info from symlink. Symlink path: \"%s\". Points to: \"%s\". Error: %s", path, realPath, err.Error())
	}
	return stat, nil
}
