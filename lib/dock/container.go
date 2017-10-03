package dock

import (
	"bytes"
	"text/template"
)

type JSONPayloadBuilder interface {
	Build() string
}

type CreateContainerPayload struct {
	Image string
	Port  string
}

func (bip CreateContainerPayload) Build() (s string, err error) {
	str := string(bip.createFromTemplateString())

	tmpl, err := template.New("").Parse(str)
	if err != nil {
		return
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(&buf, bip); err != nil {
		return
	}

	s = buf.String()

	return
}

// TODO: Left off here with formatting EXPOSED PORTS
func (bip CreateContainerPayload) createFromTemplateString() []byte {
	return []byte(
		`{{- with . -}}
		{
            "Image": "{{ .Image }}",
            "ExposedPorts": {
                {{ if .Port }}"{{ .Port }}/tcp": {}{{ else }}{}{{ end }}
            }
        }
         {{- end -}}`,
	)
}
