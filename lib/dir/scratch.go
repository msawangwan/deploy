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

// DirectoryCacher ...
type DirectoryCacher interface {
	MkTempDir(prefix string) (dir string, er error)
	FindTempDir(prefix string) (dir string, er error)
}

// DirectoryTable ...
type DirectoryTable struct {
	cache map[string]string
	sync.Mutex
}

// MkTempDir ...
func (dt *DirectoryTable) MkTempDir(prefix string) (dir string, er error) {
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
func (dt *DirectoryTable) FindTempDir(prefix string) (dir string, er error) {
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

// FlushCache ...
func FlushCache(dc DirectoryCacher) (flushcount int, er error) {
	table, ok := dc.(*DirectoryTable)
	if !ok {
		return -1, fmt.Errorf("does not implement directory cacher")
	}

	for _, dir := range table.cache {
		os.Remove(dir)
		flushcount++
	}

	return flushcount, nil
}
