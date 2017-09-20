package dock

import (
	"bytes"
	"text/template"
)

// APIStringBuilder returns an api endpoint string
type APIStringBuilder interface {
	Build() string
}

// BuildAPIURLString builds an endpoint string
func BuildAPIURLString(r APIStringBuilder) (s string, e error) {
	t, e := template.New("").Parse(r.Build())

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
