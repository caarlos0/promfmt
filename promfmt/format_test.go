package promfmt

import (
	"testing"
	"time"

	"github.com/prometheus/prometheus/promql"
	"github.com/stretchr/testify/assert"
)

func TestFormatStm(t *testing.T) {
	var assert = assert.New(t)
	out, err := AlertStmt(promql.AlertStmt{
		Name:     "my_alert",
		Duration: 1 * time.Second,
	}).Format()
	assert.NoError(err)
	assert.Equal("ALERT my_alert\n  IF <nil>\n  FOR 1s\n", out)
}

func TestCleanLabel(t *testing.T) {
	for _, data := range []struct {
		in, out string
	}{
		{
			in:  "foo from {{sdasd }} exploded",
			out: "foo from {{ sdasd }} exploded",
		},
		{
			in:  "asd {{$sdasd }} sadas",
			out: "asd {{ $sdasd }} sadas",
		},
		{
			in:  "asdasdasd {{sdasd}} sadas",
			out: "asdasdasd {{ sdasd }} sadas",
		},
		{
			in:  "asdasdasd {{ sdasd }} sadas",
			out: "asdasdasd {{ sdasd }} sadas",
		},
	} {
		t.Run(data.in, func(t *testing.T) {
			var s = cleanLabels(data.in)
			if s != data.out {
				t.Error("expected", data.out, "but got", s)
			}
		})
	}
}

func TestCleanDuration(t *testing.T) {
	for _, data := range []struct {
		d time.Duration
		s string
	}{
		{
			d: 5 * time.Hour,
			s: "5h",
		},
		{
			d: duration("1h30m"),
			s: "90m",
		},
		{
			d: duration("7h8m0s"),
			s: "428m",
		},
		{
			d: duration("0h10m0s"),
			s: "10m",
		},
	} {
		t.Run(data.s, func(t *testing.T) {
			var s = cleanDuration(data.d)
			if s != data.s {
				t.Error("expected", data.s, "but got", s)
			}
		})
	}
}

func duration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(err.Error())
	}
	return d
}
