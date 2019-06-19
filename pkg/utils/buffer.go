package utils

const defaultBufferSize int = 4 * 1024 * 1024

func CheckBufferSize(size int) int {
	if size < 512 {
		return defaultBufferSize
	}
	return size
}
