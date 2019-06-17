package threadSafe

import "sync"

// Counter is a numeric counter safe for concurrent use
type Counter struct {
	i int
	mutex sync.Mutex
}

// Increase increases the counter's value by 1
func (c *Counter) Increase() {
	c.mutex.Lock()
	c.i++
	c.mutex.Unlock()
}

// Value gets the counter's value
func (c *Counter) Value() int {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.i
}
