package dockutil

import (
	"testing"
)

type mockAPIResolver struct {
	MockEndPoint string
	MockParams   []string
}

func (m mockAPIResolver) Resolve() string {
	return `{{ with . }}
		mock endpoint: {{ .MockEndPoint }} 
		mock params: {{ range .MockParams }}
			{{ . }}
		{{ end }}
	{{ end }}`
}

func TestResolveEndpoint(t *testing.T) {
	m := &mockAPIResolver{
		"SOME ENDPOINT",
		[]string{"ONE STRING", "TWO STRING", "THREE STRING"},
	}

	result, err := ResolveAPIEndpoint(m)

	if err != nil {
		t.Errorf("%s", err)
	} else {
		t.Logf("%s", result)
	}
}
