package dock

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/msawangwan/ci.io/lib/dock"
)

// APIStringBuilder returns an api endpoint string
type APIStringBuilder interface {
	Build() []byte
}

// BuildAPIURLString builds an endpoint string
func BuildAPIURLString(r APIStringBuilder) (s string, e error) {
	helpers := template.FuncMap{
		"buildURLString": func(c, o, id string) string {
			var s string

			if id == "" {
				s = fmt.Sprintf("%s/%s", c, o)
			} else {
				s = fmt.Sprintf("%s/%s/%s", c, id, o)
			}

			return strings.TrimSuffix(s, "/")
		},
		"buildQueryString": func(q map[string]string) string {
			var s string

			for k, v := range q {
				s += fmt.Sprintf("%s=%s&", k, v)
			}

			return strings.TrimSuffix(s, "&")
		},
	}

	t, e := template.New("").Funcs(helpers).Parse(string(r.Build()))
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

// URLComponents are the individual strings that compose a docker api command url
type URLComponents struct {
	Method     string
	Command    string
	Option     string
	Parameters map[string]string
}

// ContainerCommand represents a command url: containers/{command}?<param>&<param>& etc
type ContainerCommand struct {
	URLComponents URLComponents
}

// NewContainerCommand returns a container command given the command and option
func NewContainerCommand(m, c, o string) ContainerCommand {
	return ContainerCommand{
		newURLComponents(m, c, o),
	}
}

// Build satisfies the APIEndpointResolver interface
func (c ContainerCommand) Build() []byte {
	return []byte(
		`{{- with $c := .URLComponents -}}
			{{- $e := buildURLString $c.Command $c.Option "" -}}
			{{- printf "%s" $e -}}
			{{- if $c.Parameters -}}
				{{- $q := buildQueryString $c.Parameters -}}
				{{- printf "?%s" $q -}}
			{{- end -}}
		{{- end -}}`,
	)
}

// ContainerCommandByID represents a command url: containers/{id}/{command}?<param_0>&<param_1>& etc
type ContainerCommandByID struct {
	URLComponents URLComponents
	ID            string
}

// NewContainerCommandByID returns a container command given the command and option and id
func NewContainerCommandByID(m, c, o, id string) ContainerCommandByID {
	return ContainerCommandByID{
		URLComponents: newURLComponents(m, c, o),
		ID:            id,
	}
}

// Build satisfies the APIEndpointResolver interface
func (c ContainerCommandByID) Build() []byte {
	return []byte(
		`{{- with . -}}
			{{- $e := buildURLString .URLComponents.Command .URLComponents.Option .ID -}}
			{{- printf "%s" $e -}}
			{{- if .URLComponents.Parameters -}}
				{{- $q := buildQueryString .URLComponents.Parameters -}}
				{{- printf "?%s" $q -}}
			{{- end -}}
		{{- end -}}`,
	)
}

// NewCommand returns a docker api command
// func NewCommand(c, id string) []byte {
// 	switch c {
// 	case "inspect":
// 		return NewContainerCommandByID("GET", "containers", c, id).Build()
// 	case "stop":
// 		return dock.NewContainerCommandByID("POST", "containers", c, id).Build()
// 	case "remove":
// 		return dock.NewContainerCommandByID("DELETE", "containers", "", id).Build()
// 	}

// 	return nil
// }

func newURLComponents(m, c, o string) URLComponents {
	return URLComponents{Method: m, Command: c, Option: o}
}
