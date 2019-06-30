package utils

import "os"

func IsSymLink(mode os.FileMode) bool {
	return mode & os.ModeSymlink != 0
}
