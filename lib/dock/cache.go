package dock

import (
	"fmt"
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
func (c *IDCache) Fetch(k string) (v string, e error) {
	c.Lock()
	defer c.Unlock()
	{
		if id, found := c.store[k]; found {
			v = id
		} else {
			return v, fmt.Errorf("[warn] no entry in cache for: %s", k)
		}
	}

	return
}

// Flush ...
func (c *IDCache) Flush() (n int, e error) {
	return
}

// Map ...
func (c *IDCache) Map(apply func(v string) error) error {
	var e error

	c.Lock()
	defer c.Unlock()
	for _, v := range c.store {
		e = apply(v)
		if e != nil {
			return e
		}
	}

	return nil
}

// NewIDCache ...
func NewIDCache() *IDCache {
	return &IDCache{
		store: map[string]string{},
	}
}
