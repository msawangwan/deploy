package dock

import (
	"bytes"
	"text/template"
)

type JSONPayloadBuilder interface {
	Build() (s string, err error)
}

type jsonTemplateRenderer interface {
	render() []byte
}

type CreateContainerPayload struct {
	Image string
	Port  string
}

func (ccp CreateContainerPayload) Build() (s string, err error) {
	return renderTmpl(ccp)
}

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

type StartContainerPayload struct {
	ContainerID   string
	ContainerPort string
	HostIP        string
	HostPort      string
}

func (scp StartContainerPayload) Build() (s string, err error) {
	return renderTmpl(scp)
}

func (scp StartContainerPayload) render() []byte {
	return []byte(
		`{{- with . -}}
        {
            "ID": {{ .ContainerID }},
            "PortBindings": {{ if .ContainerPort }}{
                "{{ .ContainerPort }}/tcp": {{ if is_at_least_one_not_null .HostIP .HostPort }}{
                    {{ if .HostIP }}"HostIP": "{{ .HostIP }}"{{ end }}
                    {{ if .HostPort }}"HostPort" : "{{ .HostPort }}"{{ end }}
                }{{ else }}{}{{ end }}{{ else }}{}{{ end }}
        }
         {{- end -}}`,
	)
}

func renderTmpl(jtr jsonTemplateRenderer) (s string, err error) {
	helper := template.FuncMap{
		"is_at_least_one_not_null": func(ss ...string) bool {
			for _, s := range ss {
				if s != "" {
					return true
				}
			}
			return false
		},
	}
	str := string(jtr.render())

	tmpl, err := template.New("").Funcs(helper).Parse(str)
	if err != nil {
		return
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(&buf, jtr); err != nil {
		return
	}

	s = buf.String()

	return
}
