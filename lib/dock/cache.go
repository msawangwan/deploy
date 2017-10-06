package dock

import (
	"strings"
	"sync"

	"github.com/msawangwan/ci.io/lib/cache"
)

// ImageCache ...
type ImageCache struct {
	sync.Mutex
	store map[string]string
}

// Store ...
func (c *ImageCache) Store(k, v string) {

}

// Fetch ...
func (c *ImageCache) Fetch(k string) string {
	return k
}

// Flush ...
func (c *ImageCache) Flush() (n int, e error) {
	return
}

// ContainerCache ...
type ContainerCache struct {
	sync.Mutex
	store map[string]string
}

// Store ...
func (c *ContainerCache) Store(k, v string) {

}

// Fetch ...
func (c *ContainerCache) Fetch(k string) string {
	return k
}

// Flush ...
func (c *ContainerCache) Flush() (n int, e error) {
	return
}

// NewCache ...
func NewCache(cacheType string) cache.KVStorer {
	ct := strings.ToLower(cacheType)

	switch ct {
	case "image":
		return &ImageCache{}
	case "container":
		return &ContainerCache{}
	default:
		return &cache.NullKVStore
	}
}
