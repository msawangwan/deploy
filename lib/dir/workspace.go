package dir

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// WorkspaceCache ...
type WorkspaceCache struct {
	sync.Mutex
	store map[string]string
}

// Store ...
func (wc *WorkspaceCache) Store(k, v string) {
	wc.Lock()
	defer wc.Unlock()
	{
		wc.store[k] = v
	}
}

// Fetch ...
func (wc *WorkspaceCache) Fetch(k string) (v string, e error) {
	wc.Lock()
	defer wc.Unlock()
	{
		if stored, found := wc.store[k]; found {
			v = stored
		} else {
			return v, fmt.Errorf("[warning] workspace does not exist for key: %s", k)
		}
	}

	return
}

// Flush ...
func (wc *WorkspaceCache) Flush() (n int, e error) {
	for _, v := range wc.store {
		er := os.Remove(v)
		if er != nil {
			return
		}
		n++
	}
	return
}

// NewWorkspaceCache ...
func NewWorkspaceCache() *WorkspaceCache {
	return &WorkspaceCache{
		store: map[string]string{},
	}
}

// MkTempWorkspace ...
func MkTempWorkspace(prefix string) (d string, e error) {
	d, e = ioutil.TempDir("./", prefix)
	if e != nil {
		return
	}

	return
}
