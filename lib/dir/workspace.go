package dir

import (
	"io/ioutil"
	"sync"
)

// TempDirDoesNotExistError ...
// type TempDirDoesNotExistError struct {
// 	DirPrefix string
// }

// // Error ...
// func (e TempDirDoesNotExistError) Error() string {
// 	return fmt.Sprintf("err: temp directory [%s] does not exist", e.DirPrefix)
// }

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
func (wc *WorkspaceCache) Fetch(k string) string {
	var d string

	wc.Lock()
	defer wc.Unlock()
	{
		if stored, found := wc.store[k]; found {
			d = stored
		}
	}

	return d
}

// Flush ...
func (wc *WorkspaceCache) Flush() (n int, e error) {
	return
}

// NewWorkspaceCachee ...
func NewWorkspaceCachee() *WorkspaceCache {
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

/* DEPRECATED BELOW THIS COMMENT */

// // WorkspaceCacher ...
// type WorkspaceCacher interface {
// 	MkTempDir(prefix string) (d string, er error)
// 	FindTempDir(prefix string) (d string, er error)
// 	FlushAll() (flushcount int, er error)
// }

// NewWorkspaceCache ...
// func NewWorkspaceCache() WorkspaceCacher {
// 	return &WorkspaceTable{
// 		cache: map[string]string{},
// 	}
// }

// WorkspaceTable ...
// type WorkspaceTable struct {
// 	cache map[string]string
// 	sync.Mutex
// }

// // MkTempDir ...
// func (dt *WorkspaceTable) MkTempDir(prefix string) (d string, er error) {
// 	dir, er = ioutil.TempDir("./", prefix)
// 	if er != nil {
// 		return
// 	}

// 	dt.Lock()
// 	defer dt.Unlock()
// 	{
// 		dt.cache[prefix] = dir
// 	}

// 	return
// }

// // FindTempDir ...
// func (dt *WorkspaceTable) FindTempDir(prefix string) (d string, er error) {
// 	dt.Lock()
// 	defer dt.Unlock()
// 	{
// 		if cached, found := dt.cache[prefix]; found {
// 			d = cached
// 		} else {
// 			er = TempDirDoesNotExistError{prefix}
// 		}
// 	}

// 	return
// }

// // FlushAll ...
// func (dt *WorkspaceTable) FlushAll() (flushcount int, er error) {
// 	dt.Lock()
// 	defer dt.Unlock()
// 	{
// 		if len(dt.cache) > 0 {
// 			for _, d := range dt.cache {
// 				os.Remove(d)
// 				flushcount++
// 			}
// 		}
// 	}

// 	return
// }
