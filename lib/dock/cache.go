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

}

// Flush ...
func (c *ImageCache) Flush() (n int, e error) {

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

}

// Flush ...
func (c *ContainerCache) Flush() (n int, e error) {

}

// NewCache ...
func NewCache(cacheType string) cache.KVStorer {
	ct := strings.ToLower(cacheType)

	switch ct {
	case ct == "image":
		return &ImageCache{}
	case ct == "container":
		return &ContainerCache{}
	default:
		return &cache.NullKVStore
	}
}
