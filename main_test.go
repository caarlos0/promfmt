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

	out, err := processFile("testdata/in.rules", options{
		check: false,
		diffs: false,
		write: false,
	})
	assert.NoError(err)

	expected, err := ioutil.ReadFile("testdata/out.rules")
	assert.NoError(err)
	// ioutil.WriteFile("testdata/out.rules", []byte(out), 0644)
	assert.Equal(string(expected), out)
}

func TestProcessFile(t *testing.T) {
	var assert = assert.New(t)
	before, err := ioutil.ReadFile("testdata/in.rules")
	assert.NoError(err)
	_, err = processFile("testdata/in.rules", options{
		check: true,
		diffs: true,
		write: false,
	})
	assert.Error(err)
	after, err := ioutil.ReadFile("testdata/in.rules")
	assert.NoError(err)
	assert.Equal(string(before), string(after))
}

func TestProcessAndWriteFile(t *testing.T) {
	var assert = assert.New(t)
	expected, err := ioutil.ReadFile("testdata/out.rules")
	assert.NoError(err)
	var file = filepath.Join(os.TempDir(), "test.rules")
	assert.NoError(ioutil.WriteFile(file, expected, 0644))
	_, err = processFile(file, options{
		check: false,
		diffs: false,
		write: true,
	})
	assert.NoError(err)
	after, err := ioutil.ReadFile(file)
	assert.NoError(err)
	assert.Equal(string(expected), string(after))
}

func TestFormatInvalidFile(t *testing.T) {
	var assert = assert.New(t)
	_, err := processFile("testdata/invalid.rules", options{
		check: false,
		diffs: false,
		write: false,
	})
	assert.Error(err)
}

func TestFormatFileDontExist(t *testing.T) {
	var assert = assert.New(t)
	_, err := processFile("testdata/nope.rules", options{
		check: false,
		diffs: false,
		write: false,
	})
	assert.Error(err)
}
