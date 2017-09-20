package dock

// URLComponents is a api url
type URLComponents struct {
	Command    string
	Option     string
	Parameters map[string]string
}

/*
	the two formats for docker container commands via api are:
	containers/{command}?<param_0>&<param_1>& ...
	containers/{id}/{command}?<param_0>&<param_1>& ...
*/

// ContainerCommand represents any command given to the docker host
type ContainerCommand struct {
	URLComponents URLComponents
}

// Resolve satisfies the APIEndpointResolver interface
func (c ContainerCommand) Resolve() string {
	return `{{with $c := .URLComponents}}{{$c.Command}}/{{$c.Option}}{{if $c.Parameters}}?{{range $k, $v := $c.Parameters}}{{$k}}={{$v}}&{{end}}{{end}}{{end}}`
}

// ContainerCommandByID represents any command given to a unique docker container
type ContainerCommandByID struct {
	URLComponents URLComponents
	ID            string
}

// Resolve satisfies the APIEndpointResolver interface
func (c ContainerCommandByID) Resolve() string {
	return `{{with .}}{{.URLComponents.Command}}/{{.ID}}/{{.URLComponents.Option}}{{if .URLComponents.Parameters}}?{{range $k, $v := .URLComponents.Parameters}}{{$k}}={{$v}}&{{end}}{{end}}{{end}}`
}
