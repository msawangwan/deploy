package dock

// NewCreateContainerPayload ...
func NewCreateContainerPayload(fromImg, containerPort, hostIP, hostPort string) CreateContainerPayload {
	return CreateContainerPayload{
		Image: fromImg,
		ExposedPorts: map[string]struct{}{
			containerPort: struct{}{},
		},
		OpenStdin: true,
		HostConfig: HostConfig{
			PortBindings: map[string][]HostNetworkSettings{
				containerPort: []HostNetworkSettings{
					HostNetworkSettings{
						HostIP:   hostIP,
						HostPort: hostPort,
					},
				},
			},
		},
	}
}

// CreateContainerPayload ...
type CreateContainerPayload struct {
	Image        string              `json:"Image,omitempty"`
	ExposedPorts map[string]struct{} `json:"ExposedPorts,omitempty"`
	HostConfig   HostConfig          `json:"hostConfig,omitempty"`
}

// Build ...
func (ccp CreateContainerPayload) Build() ([]byte, error) { return renderJSON(ccp) }

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
