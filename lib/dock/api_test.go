package dock

import (
	"testing"
)

type mockAPICommandURL struct {
	MockEndPoint string
	MockParams   []string
}

func (m mockAPICommandURL) Build() string {
	return `{{- with . -}}
		mock endpoint: {{- .MockEndPoint -}} 
		mock params: {{- range .MockParams -}}
			{{- . -}}&
		{{- end -}}
	{{- end -}}`
}

func TestBuildAPIURLStrings(t *testing.T) {
	m := &mockAPICommandURL{
		"SOME_ENDPOINT",
		[]string{"ONE_STRING", "TWO_STRING", "THREE_STRING"},
	}

	result, err := BuildAPIURLString(m)

	if err != nil {
		t.Errorf("%s", err)
	} else {
		t.Logf("%s", result)
	}
}
