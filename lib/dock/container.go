package dock

// NewCreateContainerPayload ...
func NewCreateContainerPayload(fromImg, containerPort string) CreateContainerPayload {
	return CreateContainerPayload{
		Image: fromImg,
		ExposedPorts: map[string]struct{}{
			containerPort: struct{}{},
		},
	}
}

// NewStartContainerPayload ...
func NewStartContainerPayload(containerID, containerPort, hostIP, hostPort string) StartContainerPayload {
	return StartContainerPayload{
		ID: containerID,
		PortBindings: map[string]HostNetworkSettings{
			containerPort: HostNetworkSettings{
				HostIP:   hostIP,
				HostPort: hostPort,
			},
		},
	}
}

// CreateContainerPayload ...
type CreateContainerPayload struct {
	Image        string              `json:"Image,omitempty"`
	ExposedPorts map[string]struct{} `json:"ExposedPorts,omitempty"`
}

// Build ...
func (ccp CreateContainerPayload) Build() ([]byte, error) { return renderJSON(ccp) }

// StartContainerPayload ...
type StartContainerPayload struct {
	ID           string                         `json:"Id,omitempty"`
	PortBindings map[string]HostNetworkSettings `json:"PortBindings,omitempty"`
}

// HostNetworkSettings ...
type HostNetworkSettings struct {
	HostIP   string `json:"HostIp,omitempty"`
	HostPort string `json:"HostPort,omitempty"`
}

// Build ...
func (scp StartContainerPayload) Build() ([]byte, error) { return renderJSON(scp) }

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
