package dock

// CreateContainerPayload ...
type CreateContainerPayload struct {
	Image string
	Port  string
}

func NewCreateContainerPayload(fromImg, containerPort string) CreateContainerPayload_new {
    return CreateContainerPayload_new {
        Image: fromImg,
        ExposedPorts: map[string]struct{}{
            containerPort: struct{},
        },
    }
}

type CreateContainerPayload_new struct {
	Image        string
	ExposedPorts map[string]struct{}
}

// Build ...
//func (ccp CreateContainerPayload) Build() ([]byte, error) { return renderTmpl(ccp) }
func (ccp CreateContainerPayload) Build() ([]byte, error) { return renderJSON(ccp) }

func (ccp CreateContainerPayload) render() string {
	return `
		{{- with . -}}
			{
				"Image": "{{ .Image }}",
				"ExposedPorts": {{ if .Port }}{
					"{{ .Port }}/tcp": {}
				}{{ else }}{}{{ end }}
			}
		{{- end -}}
		`
}

// StartContainerPayload ...
type StartContainerPayload struct {
	ContainerID   string
	ContainerPort string
	HostIP        string
	HostPort      string
}

// Build ...
func (scp StartContainerPayload) Build() ([]byte, error) { return renderJSON(scp) }

func (scp StartContainerPayload) render() string {
	return `
	{{- with . -}}
		{
			"ID": "{{ .ContainerID }}",
			"PortBindings": {{ if .ContainerPort }}{
				"{{ .ContainerPort }}/tcp": {{ if is_at_least_one_not_null .HostIP .HostPort }}{
                    {{ $c := num_elements_non_empty .HostIP .HostPort }}
                    {{ if eq $c 1 }}
					    {{ if .HostIP }}"HostIP": "{{ .HostIP }}"{{ end }}
					    {{ if .HostPort }}"HostPort": "{{ .HostPort }}"{{ end }}
                    {{ else if eq $c 2}}
                                        "HostIP": "{{ .HostIP }}",
                                        "HostPort": "{{ .HostPort }}"
                    {{ end }}
				}{{ else }}{}{{ end }}{{ else }}{}{{ end }}
		}
	{{- end -}}
	`
}

// CreateContainerAPICall ...
type CreateContainerAPICall struct {
	Parameters map[string]string
}

// Build ...
func (ccac CreateContainerAPICall) Build() ([]byte, error) { return renderTmpl(ccac) }

func (ccac CreateContainerAPICall) render() string {
	return `
	{{- with . -}}
		/containers/create
			{{- if .Parameters -}}
				?{{- append_query_parameters .Parameters -}}
			{{- end -}}
	{{- end -}}
	`
}

// StartContainerAPICall ...
type StartContainerAPICall struct {
	ContainerID string
}

// Build ...
func (scac StartContainerAPICall) Build() ([]byte, error) { return renderTmpl(scac) }

func (scac StartContainerAPICall) render() string {
	return `
	{{- with . -}}
		/containers/{{- .ContainerID -}}/start
	{{- end -}}
	`
}
