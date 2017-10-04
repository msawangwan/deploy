package dock

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
)

// Templater ...
type Templater interface {
	Build() ([]byte, error)
}

type renderer interface {
	render() string
}

func renderTmpl(r renderer) (b []byte, err error) {
	helper := template.FuncMap{
		"is_at_least_one_not_null": func(ss ...string) bool {
			for _, s := range ss {
				if s != "" {
					return true
				}
			}
			return false
		},
		"num_elements_non_empty": func(ss ...string) int {
			var count = 0

			for _, s := range ss {
				if len(s) > 0 {
					count++
				}
			}

			return count
		},
		"build_url_string": func(ss ...string) string {
			var u string

			for _, s := range ss {
				u += fmt.Sprintf("/%s", s)
			}

			return strings.TrimSuffix(u, "/")
		},
		"append_query_parameters": func(qs map[string]string) string {
			var p string

			for k, v := range qs {
				p += fmt.Sprintf("%s=%s&", k, v)
			}

			return strings.TrimSuffix(p, "&")
		},
	}

	tmpl, err := template.New("").Funcs(helper).Parse(r.render())
	if err != nil {
		return
	}

	var buf bytes.Buffer

	if err = tmpl.Execute(&buf, r); err != nil {
		return
	}

	b = buf.Bytes()

	return
}
