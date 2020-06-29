package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func printReadCloser(readCloser *io.ReadCloser, writer *os.File) {
	for {
		tmp := make([]byte, 1024)
		_, err := (*readCloser).Read(tmp)
		(*writer).Write(tmp)
		if err != nil {
			break
		}
	}
}

// RunCmdPassthrough runs a command streaming it's stdout to stdout and stderr to stderr
func RunCmdPassthrough(name string, arg ...string) error {
	cmd := exec.Command(name, arg...)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	go printReadCloser(&stdout, os.Stdout)
	go printReadCloser(&stderr, os.Stderr)
	cmd.Start()
	cmd.Wait()
	return err
}

// GetFlags gets the command flags from the environment
func GetFlags(cmd string) []string {
	flags := []string{}
	lowerCmd := strings.ToLower(cmd)
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		name := strings.ToLower(pair[0])
		value := pair[1]
		if strings.HasPrefix(name, lowerCmd) {
			name = strings.Replace(name, cmd+"_", "", 1)
			name = strings.Replace(name, "_", "-", 1)
			if len(name) <= 1 {
				os.Stderr.WriteString("flags of length 1 are not supported due to upper/lower case ambiguity\n")
				os.Exit(1)
			}
			name = "--" + name
			flags = append(flags, name)
			if value != "" {
				flags = append(flags, value)
			}
		}
	}
	return flags
}

func main() {
	cmd := os.Args[1]
	args := os.Args[2:]
	flags := GetFlags(cmd)
	args = append(flags, args...)
	fmt.Println(fmt.Sprintf("env-to-flags executing command: %s %v", cmd, strings.Join(args, " ")))
	RunCmdPassthrough(cmd, args...)
}
