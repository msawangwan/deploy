package tmpl

type TemplateData struct {
	CommandPrefix string
	CommandSuffix string
	Id            string
	QueryStrings  map[string]string
}

func (t TemplateData) Resolve() string {
	return `{{ .CommandPrefix }}/{{ if .Id }}{{ .Id}}{{ else }}{{ end }}/{{ .CommandSuffix }}`
}
