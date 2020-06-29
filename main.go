package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func printReadCloser(readCloser io.Reader, writer io.Writer) {
	for {
		tmp := make([]byte, 1024)
		_, err := readCloser.Read(tmp)
		writer.Write(tmp)
		if err != nil {
			break
		}
	}
}

func runCmdPassthroughCustomIO(name string, arg []string, stdoutWriter io.Writer, stderrWriter io.Writer) error {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go printReadCloser(stdout, stdoutWriter)
	go printReadCloser(stderr, stderrWriter)
	cmd.Start()
	cmd.Wait()
	return err
}

func getFlags(cmd string, environ []string) ([]string, error) {
	flags := []string{}
	lowerCmd := strings.ToLower(cmd)
	for _, e := range environ {
		pair := strings.SplitN(e, "=", 2)
		name := strings.ToLower(pair[0])
		value := pair[1]
		if strings.HasPrefix(name, lowerCmd) {
			name = strings.Replace(name, cmd+"_", "", 1)
			name = strings.Replace(name, "_", "-", 1)
			if len(name) <= 1 {
				return flags, errors.New("flags of length 1 are not supported due to upper/lower case ambiguity")
			}
			name = "--" + name
			flags = append(flags, name)
			if value != "" {
				flags = append(flags, value)
			}
		}
	}
	return flags, nil
}

func foo(mainArgs []string, environ []string, stdout io.Writer, stderr io.Writer) {
	cmd := mainArgs[0]
	args := mainArgs[1:]
	flags, err := getFlags(cmd, environ)
	if err != nil {
		panic(err)
	}
	args = append(flags, args...)
	fmt.Println(fmt.Sprintf("env-to-flags executing command: %s %v", cmd, strings.Join(args, " ")))
	runCmdPassthroughCustomIO(cmd, args, stdout, stderr)
}

func main() {
	foo(os.Args[1:], os.Environ(), os.Stdin, os.Stderr)
}
