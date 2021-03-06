package dir

import (
	"fmt"
	"os"
	"testing"

	"github.com/msawangwan/ci.io/lib/symbol"
)

func TestBasicWorkspaceCacherImpl(t *testing.T) {
	var dc *WorkspaceCache

	dc = NewWorkspaceCache()

	t.Logf("created a directory cacher: %s", symbol.PassMark)

	pre := "test_tmp_dir"

	dir, er := MkTempWorkspace(pre)
	if er != nil {
		t.Fatalf("%s %s", er, symbol.FailMark)
	}

	defer os.RemoveAll(dir)

	dc.Store(pre, pre)

	t.Logf("created and cached dir: %s %s", dir, symbol.PassMark)

	dir, e := dc.Fetch(pre)
	if e != nil {
		t.Fatalf("%s", e)
	}

	t.Logf("fetched dir: %s %s", dir, symbol.PassMark)
}

func TestFlushingTheCache(t *testing.T) {
	var dc *WorkspaceCache

	dc = NewWorkspaceCache()

	t.Logf("created a directory cacher: %s", symbol.PassMark)

	for i := 0; i < 20; i++ {
		dir, er := MkTempWorkspace(fmt.Sprintf("test_prefix_%d", i))
		if er != nil {
			t.Fatalf("%s", er)
		}

		dc.Store(string(i), dir)

		t.Logf("created a tmpdir: %s", dir)
	}

	cleared, er := dc.Flush()
	if er != nil {
		t.Fatalf("%s", er)
	}

	t.Logf("cache clear count: %d", cleared)

	if cleared != 20 {
		t.Fatalf("didn't delete the expected number of dirs: [expected %d][got %d] %s", 20, cleared, symbol.FailMark)
	}

	t.Logf("cache cleared: %d %s", cleared, symbol.PassMark)
}
