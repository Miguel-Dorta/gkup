package pkg

import "sync"

type safeStringList struct {
	list []string
	pos int
	mutex sync.Mutex
}

type safeFileList struct {
	list []*file
	mutex sync.Mutex
}

type safeCounter struct {
	i int
	mutex sync.Mutex
}

func makeConcurrentStringList(l []string) *safeStringList {
	return &safeStringList{list: l}
}

func makeConcurrentFileList(l []*file) *safeFileList {
	return &safeFileList{list: l}
}

func (l *safeStringList) next() *string {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.pos >= len(l.list) {
		return nil
	}

	return &l.list[l.pos]
}

func (l *safeFileList) append(f *file) {
	l.mutex.Lock()
	l.list = append(l.list, f)
	l.mutex.Unlock()
}

func (l *safeFileList) getList() []*file {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	return l.list
}

func (c *safeCounter) increase() {
	c.mutex.Lock()
	c.i++
	c.mutex.Unlock()
}

func (c *safeCounter) value() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.i
}
