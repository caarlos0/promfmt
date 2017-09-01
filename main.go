package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
	"github.com/prometheus/common/model"
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
	f, err := os.Open(*name)
	if err != nil {
		fmt.Printf("failed to open file: %s\n", *name)
		os.Exit(1)
	}
	content, err := format(f)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if *write {
		if err := ioutil.WriteFile(*name, []byte(content), 0644); err != nil {
			fmt.Printf("failed to write file: %s\n", *name)
			os.Exit(1)
		}
		return
	}
	fmt.Println(content)
}

func format(f *os.File) (string, error) {
	var result []string
	var content []string
	var reader = bufio.NewReader(f)
	for {
		line, _, rerr := reader.ReadLine()
		if rerr != nil && rerr != io.EOF {
			return "", errors.WithMessage(rerr, "failed to read file")
		}
		var s = string(line)
		if strings.HasPrefix(s, "#") {
			result = append(result, string(s))
			continue
		}
		if s == "\n" || s == "" {
			if len(content) == 0 {
				continue
			}
			stm, err := parseStm(strings.Join(content, " "))
			if err != nil {
				return "", err
			}
			result = append(result, stm)
			content = []string{}
			if rerr == io.EOF {
				break
			}
			continue
		}
		content = append(content, s)
	}
	return strings.Join(result, "\n"), nil
}

func parseStm(content string) (string, error) {
	var result []string
	stms, err := promql.ParseStmts(content)
	if err != nil {
		return "", errors.WithMessage(err, "failed to parse file")
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
			result = append(result, stm.String()+"\n")
			continue
		}
		var buff = new(bytes.Buffer)
		if err := t.Execute(buff, alert); err != nil {
			return "", errors.WithMessage(err, "failed to format")
		}
		result = append(result, buff.String())
	}
	return strings.Join(result, "\n"), nil
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
