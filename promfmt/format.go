package promfmt

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/promql"
)

const alertTemplate = `ALERT {{ .Name }}
	IF {{ .Expr }}
	FOR {{ cleanDuration .Duration}}
	{{- if .Labels }}
	LABELS {
	{{- range $key, $value := .Labels }}
		{{ $key }} = "{{ cleanLabels $value }}",
	{{- end }}
	}
	{{- end }}
	{{- if .Annotations }}
	ANNOTATIONS {
	{{- range $key, $value := .Annotations }}
		{{ $key }} = "{{ cleanLabels $value }}",
	{{- end }}
	}
	{{- end }}
`

var tmpl = template.Must(
	template.New("formatter").Funcs(
		template.FuncMap{
			"cleanDuration": cleanDuration,
			"cleanLabels":   cleanLabels,
		},
	).Parse(alertTemplate),
)

// AlertStmt is a prometheus alert statement with more methods
type AlertStmt promql.AlertStmt

// Format this alert statement
func (a AlertStmt) Format() (string, error) {
	var buff bytes.Buffer
	var err = tmpl.Execute(&buff, a)
	return buff.String(), err
}

func cleanDuration(d time.Duration) string {
	return model.Duration(d).String()
}

func cleanLabels(v interface{}) string {
	var s = fmt.Sprintf("%v", v)
	var ss []string
	for _, f := range strings.Fields(s) {
		if strings.Contains(f, "{{") {
			f = strings.Replace(f, "{{", "{{ ", -1)
		}
		if strings.Contains(f, "}}") {
			f = strings.Replace(f, "}}", " }}", -1)
		}
		ss = append(ss, strings.Join(strings.Fields(f), " "))
	}
	return strings.Join(ss, " ")
}
