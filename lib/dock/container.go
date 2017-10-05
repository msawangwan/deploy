package dock

import "fmt"

// NewCreateContainerPayload ...
func NewCreateContainerPayload(fromImg, containerPort, hostIP, hostPort string) CreateContainerPayload {
	return CreateContainerPayload{
		Image: fromImg,
		ExposedPorts: map[string]struct{}{
			containerPort: struct{}{},
		},
		AttachStdin: true,
		OpenStdin:   true,
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
	AttachStdin  bool                `json:"AttachStdin,omitempty"`
	OpenStdin    bool                `json:"OpenStdin,omitempty"`
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

func (ccac CreateContainerAPICall) Call(prefix, version string) string {
	qp := ccac.render()
	return fmt.Sprintf("%s/%s%s", prefix, version, qp)
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

func (scac StartContainerAPICall) Call(prefix, version string) string {
	return fmt.Sprintf("%s/%s/containers/%s/start", prefix, version, scac.ContainerID)
}

type StopContainerAPICall struct {
	ContainerID string
}

func (scac StopContainerAPICall) Build() ([]byte, error) { return renderTmpl(scac) }

func (scac StopContainerAPICall) render() string {
	return `
        {{- with .-}}
            /containers/{{- .ContainerID -}}/stop
        {{- end -}}
    `
}

func (scac StopContainerAPICall) Call(prefix, version string) string {
	return fmt.Sprintf("%s/%s/containers/%s/stop", prefix, version, scac.ContainerID)
}

type KillContainerAPICall struct {
	ContainerID
}

func (kcac KillContainerAPICall) Build() ([]byte, error) { return renderTmpl(kcac) }

func (kcac KillContainerAPICall) render() string {
	return `
        {{- with . -}}
            /containers/{{- .ContainerID -}}/kill
        {{- end -}}
    `
}

func (kcac KillContainerAPICall) Call(prefix, version string) string {
	return fmt.Sprintf("%s/%s/containers/%s/kill", prefix, version, kcac.ContainerID)
}

type RemoveContainerAPICall struct {
	ContainerID
}

func (rcac RemoveContainerAPICall) Build() ([]byte, error) { return renderTmpl(rcac) }

func (rcac RemoveContainerAPICall) render() string {
	return `
        {{- with . -}}
            /containers/{{- .ContainerID -}}
        {{- end -}}
    `
}

func (rcac RemoveContainerAPICall) Call(prefix, version string) string {
	return fmt.Sprintf("%s/%s/containers/%s", prefix, version, rcac.ContainerID)
}
