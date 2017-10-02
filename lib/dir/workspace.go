package dir

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
)

// TempDirDoesNotExistError ...
type TempDirDoesNotExistError struct {
	DirPrefix string
}

// Error ...
func (e TempDirDoesNotExistError) Error() string {
	return fmt.Sprintf("err: temp directory [%s] does not exist", e.DirPrefix)
}

// WorkspaceCacher ...
type WorkspaceCacher interface {
	MkTempDir(prefix string) (dir string, er error)
	FindTempDir(prefix string) (dir string, er error)
	FlushAll() (flushcount int, er error)
}

// NewWorkspaceCache ...
func NewWorkspaceCache() WorkspaceCacher {
	return &WorkspaceTable{
		cache: map[string]string{},
	}
}

// WorkspaceTable ...
type WorkspaceTable struct {
	cache map[string]string
	sync.Mutex
}

// MkTempDir ...
func (dt *WorkspaceTable) MkTempDir(prefix string) (dir string, er error) {
	dir, er = ioutil.TempDir("./", prefix)
	if er != nil {
		return
	}

	dt.Lock()
	defer dt.Unlock()
	{
		dt.cache[prefix] = dir
	}

	return
}

// FindTempDir ...
func (dt *WorkspaceTable) FindTempDir(prefix string) (dir string, er error) {
	dt.Lock()
	defer dt.Unlock()
	{
		if cached, found := dt.cache[prefix]; found {
			dir = cached
		} else {
			er = TempDirDoesNotExistError{prefix}
		}
	}

	return
}

// FlushAll ...
func (dt *WorkspaceTable) FlushAll() (flushcount int, er error) {
	dt.Lock()
	defer dt.Unlock()
	{
		if len(dt.cache) > 0 {
			for _, d := range dt.cache {
				os.Remove(d)
				flushcount++
			}
		}
	}

	return
}
