package dock

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Templater ...
type Templater interface {
	Build() (s string, err error)
}

// CreateContainerPayload ...
type CreateContainerPayload struct {
	Image string
	Port  string
}

// Build ...
func (ccp CreateContainerPayload) Build() (s string, err error) { return renderTmpl(ccp) }

func (ccp CreateContainerPayload) render() []byte {
	return []byte(
		`{{- with . -}}
		{
            "Image": "{{ .Image }}",
            "ExposedPorts": {{ if .Port }}{
                "{{ .Port }}/tcp": {}
            }{{ else }}{}{{ end }}
        }
         {{- end -}}`,
	)
}

// StartContainerPayload ...
type StartContainerPayload struct {
	ContainerID   string
	ContainerPort string
	HostIP        string
	HostPort      string
}

// Build ...
func (scp StartContainerPayload) Build() (s string, err error) { return renderTmpl(scp) }

func (scp StartContainerPayload) render() []byte {
	return []byte(
		`{{- with . -}}
        {
            "ID": {{ .ContainerID }},
            "PortBindings": {{ if .ContainerPort }}{
                "{{ .ContainerPort }}/tcp": {{ if is_at_least_one_not_null .HostIP .HostPort }}{
                    {{ if .HostIP }}"HostIP": "{{ .HostIP }}"{{ end }}
                    {{ if .HostPort }}"HostPort": "{{ .HostPort }}"{{ end }}
                }{{ else }}{}{{ end }}{{ else }}{}{{ end }}
        }
         {{- end -}}`,
	)
}

// CreateContainerAPICall ...
type CreateContainerAPICall struct {
	Parameters map[string]string
}

// Build ...
func (ccac CreateContainerAPICall) Build() (s string, err error) { return renderTmpl(ccac) }

func (ccac CreateContainerAPICall) render() []byte {
	return []byte(
		`{{- with . -}}
			/containers/create
				{{- if .Parameters -}}
					?{{- append_query_parameters .Parameters -}}
				{{- end -}}
		{{- end -}}`,
	)
}

// StartContainerAPICall ...
type StartContainerAPICall struct {
	ContainerID string
}

// Build ...
func (scac StartContainerAPICall) Build() (s string, err error) { return renderTmpl(scac) }

func (scac StartContainerAPICall) render() []byte {
	return []byte(
		`{{- with . -}}
			/containers/{{- .ContainerID -}}/start
		{{- end -}}`,
	)
}

type renderer interface {
	render() []byte
}

func renderTmpl(r renderer) (s string, err error) {
	helper := template.FuncMap{
		"is_at_least_one_not_null": func(ss ...string) bool {
			for _, s := range ss {
				if s != "" {
					return true
				}
			}
			return false
		},
		"build_url_string": func(ss ...string) string {
			var u string

			for _, s := range ss {
				u += fmt.Sprintf("/%s", s)
			}

			return strings.TrimSuffix(u, "/")
		},
		"append_query_parameters": func(qs map[string]string) string {
			var p string

			for k, v := range qs {
				p += fmt.Sprintf("%s=%s&", k, v)
			}

			return strings.TrimSuffix(p, "&")
		},
	}

	str := string(r.render())

	tmpl, err := template.New("").Funcs(helper).Parse(str)
	if err != nil {
		return
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(&buf, r); err != nil {
		return
	}

	s = buf.String()

	return
}
