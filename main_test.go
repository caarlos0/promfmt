package main

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormat(t *testing.T) {
	var assert = assert.New(t)

	expected, err := ioutil.ReadFile("testdata/out.rules")
	assert.NoError(err)
	out, err := formatFile("testdata/in.rules")
	assert.NoError(err)

	// ioutil.WriteFile("testdata/out.rules", []byte(out), 0644)
	assert.Equal(string(expected), out)
}

func TestFormatInvalidFile(t *testing.T) {
	var assert = assert.New(t)
	_, err := formatFile("testdata/invalid.rules")
	assert.Error(err)
}
