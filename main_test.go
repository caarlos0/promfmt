package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestFormat(t *testing.T) {
	in, err := os.Open("testdata/in.rules")
	if err != nil {
		t.Error(err)
	}

	expected, err := ioutil.ReadFile("testdata/out.rules")
	if err != nil {
		t.Error(err)
	}

	out, err := format(in)
	if err != nil {
		t.Error(err)
	}
	// ioutil.WriteFile("testdata/out.rules", []byte(out), 0644)
	if string(expected) != out {
		t.Error("failed to format file, got:", out)
	}
}

func TestFormatInvalidFile(t *testing.T) {
	in, err := os.Open("testdata/invalid.rules")
	if err != nil {
		t.Error(err)
	}
	_, err = format(in)
	if err == nil {
		t.Error("expected an error, got nil")
	}
}
