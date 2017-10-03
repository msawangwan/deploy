package dock

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
