// Package tmp is intended to be used as a temporal develop directory.
// All it's content will not be here for the post-alpha versions.
package tmp

var (
	BufferSize = 4*1024*1024
	CopyBuf = make([]byte, BufferSize)
	ReadSymLink = false
)
