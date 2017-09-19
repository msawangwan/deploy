package dockutil

import (
	"bytes"
	"text/template"
)

// APIResolver returns an api endpoint string
type APIResolver interface {
	Resolve() string
}

// ResolveEndpoint builds a an endpoint string
func ResolveEndpoint(r APIResolver) (s string, e error) {
	t, e := template.New("").Parse(r.Resolve())

	if e != nil {
		return
	}

	var b bytes.Buffer

	if e = t.Execute(&b, r); e != nil {
		return
	}

	s = b.String()

	return
}
