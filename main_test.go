package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
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

func TestProcessFile(t *testing.T) {
	var assert = assert.New(t)
	before, err := ioutil.ReadFile("testdata/in.rules")
	assert.NoError(err)
	assert.NoError(processFile("testdata/in.rules", false))
	after, err := ioutil.ReadFile("testdata/in.rules")
	assert.NoError(err)
	assert.Equal(before, after)
}

func TestProcessInvalidFile(t *testing.T) {
	var assert = assert.New(t)
	assert.Error(processFile("testdata/invalid.rules", false))
}

func TestProcessAndWriteFile(t *testing.T) {
	var assert = assert.New(t)
	expected, err := ioutil.ReadFile("testdata/out.rules")
	assert.NoError(err)
	var file = filepath.Join(os.TempDir(), "test.rules")
	assert.NoError(ioutil.WriteFile(file, expected, 0644))
	assert.NoError(processFile(file, true))
	after, err := ioutil.ReadFile(file)
	assert.NoError(err)
	assert.Equal(expected, after)
}

func TestFormatInvalidFile(t *testing.T) {
	var assert = assert.New(t)
	_, err := formatFile("testdata/invalid.rules")
	assert.Error(err)
}

func TestFormatFileDontExist(t *testing.T) {
	var assert = assert.New(t)
	_, err := formatFile("testdata/nope.rules")
	assert.Error(err)
}
