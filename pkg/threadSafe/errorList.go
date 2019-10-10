package threadSafe

import "sync"

type ErrorList struct {
	list []error
	m sync.Mutex
}

func NewErrorList(e []error) *ErrorList {
	if e == nil {
		e = make([]error, 0, 100)
	}
	return &ErrorList{list:e}
}

func (el *ErrorList) Append(err error) {
	el.m.Lock()
	el.list = append(el.list, err)
	el.m.Unlock()
}

func (el *ErrorList) GetList() []error {
	el.m.Lock()
	defer el.m.Unlock()
	return el.list
}
