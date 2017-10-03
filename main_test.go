package main

// func TestCachesNonNil(t *testing.T) {
// 	dirCache.Lock()
// 	defer dirCache.Unlock()
// 	{
// 		dirCache.store["some_dir"] = "dirpath"
// 		t.Logf("%+v", dirCache)
// 	}
// }

// func TestCreateTmpDirAndRemove(t *testing.T) {
// 	c := newCache()

// 	ws, e := createTmpWorkspace(c, "testdir")
// 	if e != nil {
// 		t.Fatal(e)
// 	}

// 	defer os.RemoveAll(ws)

// 	t.Logf("created tmp workspace: %s", ws)
// }

// func TestFetchTmpDirAndRemove(t *testing.T) {
// 	var testDir = "some_test_dir"

// 	c := newCache()
// 	c.store[testDir] = "some_test_dir_name"

// 	ws, e := getWorkspace(c, testDir)
// 	if e != nil {
// 		t.Fatal(e)
// 	}

// 	t.Logf("fetched tmp ws: %s", ws)
// }
