package github

import (
	"bytes"
	"testing"
)

var mockResponse = []byte(`{ "Repository": { "Name": "repo_name" } }`)

func TestExtractRepoName(t *testing.T) {
	r := bytes.NewReader(mockResponse)

	n, e := ExtractRepositoryName(r)
	if e != nil {
		t.Fatal(e)
	}

	t.Logf("%s", n)
}
