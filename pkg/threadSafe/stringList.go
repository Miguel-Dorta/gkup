package threadSafe

import "sync"

// FileList is a list of strings safe for concurrent use
type StringList struct {
	list  []string
	pos   int
	mutex sync.Mutex
}

// NewStringList creates a new StringList object
func NewStringList(l []string) *StringList {
	if l == nil {
		l = make([]string, 0, 100)
	}
	return &StringList{list: l}
}

// Next gets the reference to the next string when reading concurrently.
// Returns nil when the end of the slice is reached
func (l *StringList) Next() *string {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.pos >= len(l.list) {
		return nil
	}
	s := &l.list[l.pos]
	l.pos++

	return s
}

func (l *StringList) GetPosUnsafe() int {
	return l.pos
}

func (l *StringList) GetLenUnsafe() int {
	return len(l.list)
}
