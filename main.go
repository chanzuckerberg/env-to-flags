package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

// printReadCloser prints from a ReadCloser until it's closed
func printReadCloser(readCloser io.ReadCloser, writer io.Writer) {
	for {
		tmp := make([]byte, 1024)
		_, err := readCloser.Read(tmp)
		writer.Write(tmp)
		if err != nil {
			break
		}
	}
}

// runCmdPassthrough runs a command streaming output to provided stdout and stderr writers
func runCmdPassthrough(name string, arg []string, stdoutWriter io.Writer, stderrWriter io.Writer) error {
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
	err = cmd.Wait()
	return err
}

// getFlags parses flags for an executable given an environment
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

// runWithEnvFlags runs a command with flags parsed from environ
func runWithEnvFlags(mainArgs []string, environ []string, stdout io.Writer, stderr io.Writer) error {
	cmd := mainArgs[0]
	args := mainArgs[1:]
	flags, err := getFlags(cmd, environ)
	if err != nil {
		panic(err)
	}
	args = append(flags, args...)
	message := fmt.Sprintf("env-to-flags executing command: %s %v\n", cmd, strings.Join(args, " "))
	stderr.Write([]byte(message))
	return runCmdPassthrough(cmd, args, stdout, stderr)
}

func main() {
	err := runWithEnvFlags(os.Args[1:], os.Environ(), os.Stdin, os.Stderr)
	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
