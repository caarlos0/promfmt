package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/prometheus/prometheus/promql"
)

var write = flag.Bool("w", false, "override the source file with the formatted file")
var name = flag.String("f", "", "file to format")

func main() {
	flag.Parse()
	if *name == "" {
		fmt.Println("missing file name")
		os.Exit(2)
	}
	f, err := ioutil.ReadFile(*name)
	if err != nil {
		fmt.Printf("failed to open file: %s\n", *name)
		os.Exit(1)
	}
	var content = format(string(f))
	if *write {
		if err := ioutil.WriteFile(*name, []byte(content), 0644); err != nil {
			fmt.Printf("failed to write file: %s\n", *name)
			os.Exit(1)
		}
		return
	}
	fmt.Println(content)
}

func format(content string) string {
	var result []string
	stms, err := promql.ParseStmts(content)
	if err != nil {
		log.Fatal(err.Error())
	}

	var t = template.Must(
		template.New("formatter").Funcs(
			template.FuncMap{
				"cleanDuration": cleanDuration,
				"cleanLabels":   cleanLabels,
			},
		).Parse(alertTemplate),
	)
	for _, stm := range stms {
		alert, ok := stm.(*promql.AlertStmt)
		if !ok {
			result = append(result, stm.String())
			continue
		}
		var buff = new(bytes.Buffer)
		if err := t.Execute(buff, alert); err != nil {
			log.Fatal(err.Error())
		}
		result = append(result, buff.String())
	}
	return strings.Join(result, "\n")
}

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
		ss = append(ss, strings.Join(strings.Fields(f), " "))
	}
	return strings.Join(ss, " ")
}

var alertTemplate = `ALERT {{ .Name }}
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
