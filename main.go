package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"text/template"
	"time"

	"github.com/prometheus/prometheus/promql"
)

var tmpl = `ALERT {{ .Name }}
	IF {{ .Expr }}
	FOR {{ cleanDuration .Duration}}
	LABELS {
	{{- range $key, $value := .Labels }}
		{{ $key }} = "{{ cleanLabels $value }}",
	{{- end }}
	}
	ANNOTATIONS {
	{{- range $key, $value := .Annotations }}
		{{ $key }} = "{{ cleanLabels $value }}",
	{{- end }}
	}
`

func cleanDuration(d time.Duration) string {
	return strings.Replace(
		strings.Replace(
			strings.Replace(
				d.String(),
				"h0m", "h", 1,
			),
			"m0s", "m", 1,
		),
		"h0s", "h", 1,
	)
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
		fmt.Println(strings.Fields(f))
		ss = append(ss, strings.Join(strings.Fields(f), " "))
	}
	return strings.Join(ss, " ")
}

func main() {
	var content = `alert moises
	If a >1
	fOr 5h
	labels {
		a = "b",
	}
	ANNOTATIONS {
		SUMMARY = "{{$labels.instance}}: High memory usage detected",
		DESCRIPTION = "{{$labels.instance}}: Memory usage is high (current value is: {{ $value }})"
	  }
	`
	stms, err := promql.ParseStmts(content)
	if err != nil {
		log.Fatal(err.Error())
	}

	var t = template.Must(template.New("alert").Funcs(
		template.FuncMap{
			"cleanDuration": cleanDuration,
			"cleanLabels":   cleanLabels,
		},
	).Parse(tmpl))
	for _, stm := range stms {
		alert, ok := stm.(*promql.AlertStmt)
		if !ok {
			continue
		}
		var buff = new(bytes.Buffer)
		if err := t.Execute(buff, alert); err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(buff.String())
	}
}
