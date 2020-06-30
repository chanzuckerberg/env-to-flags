package main

import (
	"strings"
	"testing"
)

// dummyWriter is a struct that writes to an underlying string so it can be inspected
type dummyWriter struct {
	str string
}

// Write writes to the underlying string of a dummyWriter
func (dw *dummyWriter) Write(p []byte) (int, error) {
	n := 0
	for _, c := range p {
		if c != 0 {
			(*dw).str = (*dw).str + string(c)
			n++
		}
	}
	return n, nil
}

// newDummyWriter creates a new empty dummyWriter
func newDummyWriter() dummyWriter {
	return dummyWriter{str: ""}
}

func assertEqualSlices(t *testing.T, a []string, b []string) {
	if len(a) != len(b) {
		t.Errorf("Expected %v to equal %v", a, b)
		return
	}

	for i := range a {
		if a[i] != b[i] {
			t.Errorf("Expected %v to equal %v", a, b)
			return
		}
	}
}

func TestGetFlagsAllCases(t *testing.T) {
	res, err := getFlags("foo", []string{"FOO_BAR=a", "FOO_BAZ=", "BLAH_BLAH=c"})
	expected := []string{"--bar", "a", "--baz"}

	if err != nil {
		t.Errorf("Expected nil error but it was %s", err.Error())
	}

	assertEqualSlices(t, res, expected)
}

func TestGetFlagsSingleLetterError(t *testing.T) {
	_, err := getFlags("foo", []string{"FOO_BAR=a", "FOO_B=a", "FOO_BAZ=b"})

	if err == nil {
		t.Errorf("Expected getFlags to return an error on a flag of length 1 but it was nil")
		return
	}

	if err.Error() != "flags of length 1 are not supported due to upper/lower case ambiguity" {
		t.Errorf("Expected getFlags to return a length 1 error but the error was %s", err.Error())
	}
}

func TestRunCmdPassthrough(t *testing.T) {
	stdout := newDummyWriter()
	stderr := newDummyWriter()
	runCmdPassthrough("echo", []string{"foo"}, &stdout, &stderr)
	expectedStdout := "foo\n"
	expectedStderr := ""
	if stdout.str != expectedStdout {
		t.Errorf("Expected runCmdPassthrough to use echo to write '%s' to stdout but it wrote %s", expectedStdout, stdout.str)
	}

	if stderr.str != expectedStderr {
		t.Errorf("Expected runCmdPassthrough to use echo to write '%s' to stderr but it wrote %s", expectedStderr, stderr.str)
	}
}

func TestRunWithEnvFlags(t *testing.T) {
	environ := []string{"CAT_VERSION="}
	stdout := newDummyWriter()
	stderr := newDummyWriter()
	expectedStderr := "env-to-flags executing command: cat --version\n"
	runWithEnvFlags([]string{"cat"}, environ, &stdout, &stderr)
	if !strings.HasPrefix(stdout.str, "cat (GNU coreutils)") {
		t.Errorf("Expected cat to run with the version option but output was %s", stdout.str)
	}

	if stderr.str != expectedStderr {
		t.Errorf("Expected to output %s to stderror but it was %s", expectedStderr, stderr.str)
	}
}

