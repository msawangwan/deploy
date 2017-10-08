package dock

import (
	"fmt"
)

// ListImageAPICall ...
type ListImageAPICall struct {
	Parameters map[string]string
}

func (liac ListImageAPICall) render() string {
	return `
	{{- with . -}}
		/images/json
			{{- if .Parameters -}}
				?{{- append_query_parameters .Parameters -}}
			{{- end -}}
	{{- end -}}
	`
}

// Build ...
func (liac ListImageAPICall) Build() ([]byte, error) { return renderTmpl(liac) }

// BuildImageAPICall ...
type BuildImageAPICall struct {
	Parameters map[string]string
}

func (biac BuildImageAPICall) render() string {
	return `
	{{- with . -}}
		/build
			{{- if .Parameters -}}
				?{{- append_query_parameters .Parameters -}}
			{{- end -}}
	{{- end -}}
	`
}

// Build ...
func (biac BuildImageAPICall) Build() ([]byte, error) { return renderTmpl(biac) }

// Call ...
func (biac BuildImageAPICall) Call(prefix, version string) string {
	return fmt.Sprintf("%s/%s/build/%s", prefix, version, "")
}

// InspectImageAPICall ...
type InspectImageAPICall struct {
	Name string
}

func (iiac InspectImageAPICall) render() string {
	return `
	{{- with . -}}
		/images/{{ .Name }}/json
	{{- end -}}
	`
}

// Build ...
func (iiac InspectImageAPICall) Build() ([]byte, error) { return renderTmpl(iiac) }

// RemoveImageAPICall ...
type RemoveImageAPICall struct {
	Name string
}

func (riac RemoveImageAPICall) render() string {
	return `
	{{- with . -}}
		/images/{{- .Name -}}
	{{- end -}}
	`
}

// Build ...
func (riac RemoveImageAPICall) Build() ([]byte, error) { return renderTmpl(riac) }

// DeleteUnusedImagesAPICall deletes unused images, valid filters include:
// dangling=<boolean>
// until=<string>
// label=<key>=<value>
type DeleteUnusedImagesAPICall struct {
	Filters map[string]string
}

func (duiac DeleteUnusedImagesAPICall) render() string {
	return `
	{{- with . -}}
		/images/prune
			{{- if .Filters -}}
				?{{- append_query_parameters .Filters -}}
			{{- end -}}
	{{- end -}}
	`
}

// Build ...
func (duiac DeleteUnusedImagesAPICall) Build() ([]byte, error) { return renderTmpl(duiac) }
