package dock

// URLComponents are the individual strings that compose a docker api command url
type URLComponents struct {
	Command    string
	Option     string
	Parameters map[string]string
}

// ContainerCommand represents any command given to the docker host of the form:
// containers/{command}?<param>&<param>& etc
type ContainerCommand struct {
	URLComponents URLComponents
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

// ContainerCommandByID represents any command given to a unique docker container of the form:
// containers/{id}/{command}?<param_0>&<param_1>& etc
type ContainerCommandByID struct {
	URLComponents URLComponents
	ID            string
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
