package main

import (
	"strings"
	"testing"

	"golang.org/x/exp/errors/fmt"
)

func TestGetFlags(t *testing.T) {
	res, err := getFlags("foo", []string{"FOO_BAR=a", "FOO_BAZ=", "BLAH_BLAH=c"})
	expected := []string{"--bar", "a", "--baz"}

	if err != nil {
		t.Errorf("Expected nil error but it was %s", err.Error())
	}

	for i := range res {
		if res[i] != expected[i] {
			t.Errorf("Expected %s to equal %s", res[i], expected[i])
		}
	}
}

type dummyWriter struct {
	str string
}

func (dw *dummyWriter) Write(p []byte) (int, error) {
	(*dw).str = (*dw).str + string(p[:])
	return len(p), nil
}

func TestRunCmdPassthrough(t *testing.T) {
	stdout := dummyWriter{str: ""}
	stderr := dummyWriter{str: ""}
	runCmdPassthroughCustomIO("echo", []string{"foo"}, &stdout, &stderr)
	if stdout.str != "foo" {
		fmt.Errorf("asdasd")
	}
}

func TestMain(t *testing.T) {
	environ := []string{"GREP_VERSION="}
	stdout := dummyWriter{str: ""}
	stderr := dummyWriter{str: ""}
	foo([]string{"cat"}, environ, &stdout, &stderr)
	if !strings.HasPrefix(stdout.str, "cat (GNU coreutils)") {
		fmt.Errorf("Expected asdasd")
	}
}
