package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	promfmt "github.com/caarlos0/promfmt/format"
	"github.com/pkg/errors"
	"github.com/pmezard/go-difflib/difflib"
	"github.com/prometheus/prometheus/promql"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	version = "master"
	app     = kingpin.New("promfmt", "promfmt formats prometheus' .rules files")
	write   = app.Flag("write", "Override the source file with the formatted file").Short('w').Bool()
	fail    = app.Flag("fail", "Fails if the file is not in the expected format").Short('f').Bool()
	diffs   = app.Flag("diff", "Prints the diff between the file and the formatted file instead of the end result").Short('d').Bool()
	name    = app.Arg("file", "Path to file to be formatted").Required().ExistingFile()
)

type options struct {
	write, fail, diffs bool
}

func main() {
	app.Version("promfmt version " + version)
	app.HelpFlag.Short('h')
	app.VersionFlag.Short('v')
	kingpin.MustParse(app.Parse(os.Args[1:]))
	opts := options{
		write: *write,
		fail:  *fail,
		diffs: *diffs,
	}
	if _, err := processFile(*name, opts); err != nil {
		kingpin.Fatalf("%s: %v\n", *name, err)
	}
}

func processFile(name string, opts options) (string, error) {
	original, err := ioutil.ReadFile(name)
	if err != nil {
		return "", errors.Wrap(err, "failed to open file")
	}
	content, err := format(bytes.NewBuffer(original))
	if err != nil {
		return "", err
	}
	if string(original) == content {
		return content, nil
	}
	if opts.diffs {
		diff, err := diffContents(name, string(original), content)
		if err != nil {
			return content, err
		}
		fmt.Println(diff)
	}
	if opts.write {
		return content, ioutil.WriteFile(name, []byte(content), 0644)
	}
	if opts.fail {
		return content, fmt.Errorf("file does not match")
	}
	return content, nil
}

func diffContents(name, original, formatted string) (string, error) {
	return difflib.GetUnifiedDiffString(difflib.UnifiedDiff{
		A:        difflib.SplitLines(original),
		B:        difflib.SplitLines(formatted),
		FromFile: name,
		ToFile:   fmt.Sprintf("formatted %s", name),
	})
}

// TODO: since we are already reading the entire file before here, maybe
// simplify to a for-range loop on the lines or something like that...
func format(f io.Reader) (string, error) {
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
			result = append(result, s)
			continue
		}
		if s != "\n" && s != "" {
			content = append(content, s)
			continue
		}
		if len(content) == 0 {
			if eof {
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
