package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/caarlos0/promfmt/format"
	"github.com/pkg/errors"
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
	content, err := formatFile(f)
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

func formatFile(f *os.File) (string, error) {
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
		if s == "\n" || s == "" {
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
			if eof {
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
	for _, stm := range stms {
		alert, ok := stm.(*promql.AlertStmt)
		if !ok {
			result = append(result, stm.String()+"\n")
			continue
		}
		str, err := format.AlertStmt(*alert).Format()
		if err != nil {
			return "", err
		}
		result = append(result, str)
	}
	return strings.Join(result, "\n"), nil
}
