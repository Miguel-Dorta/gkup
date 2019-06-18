package threadSafe

import "sync"

// Fuse is a boolean that can only be set to true safe for concurrent use
type Fuse struct {
	value bool
	mutex sync.RWMutex
}

// Trigger sets the Fuse value to true
func (f *Fuse) Trigger() {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	if f.value {
		return
	}

	f.mutex.Lock()
	f.value = true
	f.mutex.Unlock()
}

// Value returns the value of the Fuse.
// Its value will be true if Trigger() was called one or more times, false otherwise.
func (f *Fuse) Value() bool {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	return f.value
}
