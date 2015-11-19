package kubectl

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Path is the path of the kubectl binary
var Path = "kubectl"

// Runner is an interface to wrap kubectl convenience methods
type Runner interface {
	ClusterInfo() ([]byte, error)
	Create([]byte, string, bool) ([]byte, error)
	Delete(string, string, string, bool) ([]byte, error)
	Get([]byte, string, bool) ([]byte, error)
}

// RealRunner implements Runner to execute kubectl commands
type RealRunner struct{}

// Client stores the instance of Runner
var Client Runner = RealRunner{}

func assignStdin(cmd *exec.Cmd, in []byte) {
	cmd.Stdin = bytes.NewBuffer(in)
}

func commandToString(cmd *exec.Cmd) string {
	var stdin string

	if cmd.Stdin != nil {
		buf := new(bytes.Buffer)
		buf.ReadFrom(cmd.Stdin)
		stdin = fmt.Sprintf("< %s", buf.String())
	}

	return fmt.Sprintf("%s %s", strings.Join(cmd.Args, " "), stdin)
}
