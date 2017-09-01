package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/caarlos0/promfmt/promfmt"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/promql"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	write = kingpin.Flag(
		"write", "override the source file with the formatted file",
	).Short('w').Bool()
	name = kingpin.Arg(
		"file", "path to file to be formatted",
	).Required().String()
)

func main() {
	kingpin.Parse()
	content, err := formatFile(*name)
	if err != nil {
		kingpin.Fatalf("%s: %v\n", *name, err)
	}
	if *write {
		if err := ioutil.WriteFile(*name, []byte(content), 0644); err != nil {
			kingpin.Fatalf("%s: %v\n", *name, err)
		}
		return
	}
	fmt.Println(content)
}

func formatFile(name string) (string, error) {
	f, err := os.Open(name)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file")
	}
	return format(f)
}

func format(f *os.File) (string, error) {
	var result []string
	var content []string
	var reader = bufio.NewReader(f)
	for {
		line, _, rerr := reader.ReadLine()
		var eof = rerr == io.EOF
		if rerr != nil && !eof {
			return "", errors.WithMessage(rerr, "failed to read file")
		}
		var s = string(line)
		if strings.HasPrefix(s, "#") {
			result = append(result, string(s))
			continue
		}
		if s != "\n" && s != "" {
			content = append(content, s)
			continue
		}
		if len(content) == 0 {
			if eof {
				result = append(result, "")
				break
			}
			continue
		}
		stm, err := parseStm(strings.Join(content, " "))
		if err != nil {
			return "", err
		}
		result = append(result, stm)
		content = []string{}
	}
	return strings.Join(result, "\n"), nil
}

func parseStm(content string) (string, error) {
	var result []string
	stms, err := promql.ParseStmts(content)
	if err != nil {
		return "", errors.WithMessage(err, "failed to parse file")
	}
	for _, stm := range stms {
		alert, ok := stm.(*promql.AlertStmt)
		if !ok {
			result = append(result, stm.String()+"\n")
			continue
		}
		str, err := promfmt.AlertStmt(*alert).Format()
		if err != nil {
			return "", err
		}
		result = append(result, str)
	}
	return strings.Join(result, "\n"), nil
}
