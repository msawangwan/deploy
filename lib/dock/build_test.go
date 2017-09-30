package dock

import "testing"

func TestBuildRepository(t *testing.T) {
	out, er := BuildRepository("bare_dir", "tmp_dir", "test/repo", "user", "oath")
	if er != nil {
		t.Errorf("%s", out)
		t.Fatalf("%s", er)
	}

	t.Logf("%s", out)
}
