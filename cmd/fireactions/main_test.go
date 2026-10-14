package main

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunReportsFailureWithNonZeroExit(t *testing.T) {
	// A command that cannot do its job must exit non-zero: otherwise a supervisor (or a
	// CI gate) reads the failure as success. A config file that does not exist is the
	// smallest reliable failure.
	stderr := &bytes.Buffer{}
	code := run([]string{"server", "--config", "testdata/missing-config.yaml"}, stderr)

	assert.Equal(t, 1, code)
	assert.Contains(t, stderr.String(), "Error executing command")
}

func TestRunSucceedsWithZeroExit(t *testing.T) {
	stderr := &bytes.Buffer{}
	code := run([]string{"version"}, stderr)

	assert.Equal(t, 0, code)
	assert.Empty(t, stderr.String())
}
