package utils

import (
	"fmt"
	"syscall"
)

// isHidden returns whether the path provided is hidden
func isHidden(path, name string) (bool, error) {
	ptr, err := syscall.UTF16PtrFromString(path)
	if err != nil {
		return false, fmt.Errorf("cannot get file pointer (windows syscall): %s", err.Error())
	}

	attr, err := syscall.GetFileAttributes(ptr)
	if err != nil {
		return false, fmt.Errorf("cannot get file attributes (windows syscall): %s", err.Error())
	}

	return attr&syscall.FILE_ATTRIBUTE_HIDDEN != 0, nil
}
