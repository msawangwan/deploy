package dock

// URLComponents are the individual strings that compose a docker api command url
type URLComponents struct {
	Method     string
	Command    string
	Option     string
	Parameters map[string]string
}

func newURLComponents(m, c, o string) URLComponents {
	return URLComponents{Method: m, Command: c, Option: o}
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
func (c ContainerCommand) Build() string {
	return `{{- with $c := .URLComponents -}}
				{{- $c.Command -}}/{{- $c.Option -}}
				{{- if $c.Parameters -}}
					?{{- range $k, $v := $c.Parameters -}}
						{{- $k -}}={{- $v -}}&
					{{- end -}}
				{{- end -}}
			{{- end -}}`
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
func (c ContainerCommandByID) Build() string {
	return `{{- with . -}}
				{{- .URLComponents.Command -}}/{{- .ID -}}/{{- .URLComponents.Option -}}
				{{- if .URLComponents.Parameters -}}
					?{{- range $k, $v := .URLComponents.Parameters -}}
						{{- $k -}}={{- $v -}}&
					{{- end -}}
				{{- end -}}
			{{- end -}}`
}
