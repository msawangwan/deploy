package dock

// APIURL is a api url
type APIURL struct {
	ID         string
	Command    string
	Option     string
	Parameters map[string]string
}

// List is eqivualnt to docker containers ls
type List struct {
	URLComponents APIURL
}

// Resolve satisfies the APIEndpointResolver interface
func (l List) Resolve() string {
	return `{{with $l := .URLComponents}}{{$l.Command}}/{{$l.Option}}`
}

// Create is the create command
type Create struct {
	URLComponents APIURL
}

// Resolve satisfies the APIEndpointResolver interface
func (c Create) Resolve() string {
	return `{{with $c := .URLComponents}}{{$c.Command}}/{{$c.Option}}?{{range $k, $v := $c.Parameters}}{{$k}}={{$v}}{{end}}{{end}}`
	// return `{{with .}}{{.URLComponents.Command}}/{{.URLComponents.Option}}?{{range $k, $v := .URLComponents.Parameters}}{{$k}}={{$v}}{{end}}{{end}}`
}

// Inspect is a inspect command
type Inspect struct {
	URLComponents APIURL
}
