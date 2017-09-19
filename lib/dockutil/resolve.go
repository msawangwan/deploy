package dockutil

import (
	"bytes"
	"text/template"
)

// APIEndpointResolver returns an api endpoint string
type APIEndpointResolver interface {
	Resolve() string
}

// ResolveAPIEndpoint builds an endpoint string
func ResolveAPIEndpoint(r APIEndpointResolver) (s string, e error) {
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
