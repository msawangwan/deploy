package main

import "testing"

func TestCachesNonNil(t *testing.T) {
	dirCache.Lock()
	defer dirCache.Unlock()
	{
		dirCache.store["some_dir"] = "dirpath"
		t.Logf("%+v", dirCache)
	}
}
