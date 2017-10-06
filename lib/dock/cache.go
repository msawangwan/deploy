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

// // ImageCache ...
// type ImageCache struct {
// 	sync.Mutex
// 	store map[string]string
// }

// // Store ...
// func (c *ImageCache) Store(k, v string) {

// }

// // Fetch ...
// func (c *ImageCache) Fetch(k string) string {
// 	return k
// }

// // Flush ...
// func (c *ImageCache) Flush() (n int, e error) {
// 	return
// }

// // ContainerCache ...
// type ContainerCache struct {
// 	sync.Mutex
// 	store map[string]string
// }

// // Store ...
// func (c *ContainerCache) Store(k, v string) {

// }

// // Fetch ...
// func (c *ContainerCache) Fetch(k string) string {
// 	return k
// }

// // Flush ...
// func (c *ContainerCache) Flush() (n int, e error) {
// 	return
// }

// // NewCache ...
// func NewCache(cacheType string) cache.KVStorer {
// 	ct := strings.ToLower(cacheType)

// 	switch ct {
// 	case "image":
// 		return &ImageCache{}
// 	case "container":
// 		return &ContainerCache{}
// 	default:
// 		return &cache.NullKVStore
// 	}
// }
