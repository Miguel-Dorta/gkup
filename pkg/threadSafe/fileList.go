package threadSafe

import (
	"github.com/Miguel-Dorta/gkup/pkg/files"
	"sync"
)

// FileList is a list of files.File safe for concurrent use
type FileList struct {
	list  []*files.File
	pos   int
	mutex sync.Mutex
}

// NewFileList creates a new FileList object
func NewFileList(l []*files.File) *FileList {
	return &FileList{list: l}
}

// Append appends a files.File object to the end of the list
func (l *FileList) Append(f *files.File) {
	l.mutex.Lock()
	l.list = append(l.list, f)
	l.mutex.Unlock()
}

// GetList gets the internal slice of files.File objects
func (l *FileList) GetList() []*files.File {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.list
}

// Next gets the next files.File object when reading concurrently.
// Returns nil when the end of the slice is reached
func (l *FileList) Next() *files.File {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.pos >= len(l.list) {
		return nil
	}
	f := l.list[l.pos]
	l.pos++

	return f
}
