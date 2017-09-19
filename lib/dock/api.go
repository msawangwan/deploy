package dock

// APIURL is a api url
type APIURL struct {
	ID         string
	Command    string
	Option     string
	Parameters map[string]string
}

// ListContainers is eqivualnt to docker containers ls
type ListContainers struct {
	URLComponents APIURL
}

// Resolve satisfies the APIEndpointResolver interface
func (l ListContainers) Resolve() string {
	return `{{with $l := .URLComponents}}{{$l.Command}}/{{$l.Option}}{{end}}`
}

// CreateContainer is the create command
type CreateContainer struct {
	URLComponents APIURL
}

// Resolve satisfies the APIEndpointResolver interface
func (c CreateContainer) Resolve() string {
	return `{{with $c := .URLComponents}}{{$c.Command}}/{{$c.Option}}?{{range $k, $v := $c.Parameters}}{{$k}}={{$v}}{{end}}{{end}}`
}

// InspectContainer is a inspect command
type InspectContainer struct {
	URLComponents APIURL
}

// Resolve satisfies the APIEndpointResolver interface
func (i InspectContainer) Resolve() string {
	return `{{with $i := .URLComponents}}{{$i.Command}}/{{$i.ID}}/{{$i.Option}}{{end}}`
}

// StartContainer is the docker start container command
type StartContainer struct {
	URLComponents APIURL
}

// Resolve satisfies the APIEndpointResolver interface
func (s StartContainer) Resolve() string {
	return `{{with $s := .URLComponents}}{{end}}`
}
