package dock

import (
	"sync"
)

// IDCache ...
type IDCache struct {
	sync.Mutex
	store map[string]string
}

// Store ...
func (c *IDCache) Store(k, v string) {
	c.Lock()
	defer c.Unlock()
	{
		c.store[k] = v
	}
}

// Fetch ...
func (c *IDCache) Fetch(k string) string {
	var id string

	c.Lock()
	defer c.Unlock()
	{
		if v, found := c.store[k]; found {
			id = v
		}
	}

	return id
}

// Flush ...
func (c *IDCache) Flush() (n int, e error) {
	return
}

// NewIDCache ...
func NewIDCache() *IDCache {
	return &IDCache{
		store: map[string]string{},
	}
}
