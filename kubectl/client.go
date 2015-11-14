package kubectl

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

var (
	Path   string    = "kubectl"
	Writer io.Writer = os.Stderr
)

func kubectlCmd(args ...string) *exec.Cmd {
	cmd := exec.Command(Path, args...)
	return cmd
}

// kubectlBuilder is used to build, custimize and execute a kubectl Command.
type kubectlBuilder struct {
	cmd    *exec.Cmd
	dryRun bool
}

func (b kubectlBuilder) DryRun(dryRun bool) *kubectlBuilder {
	b.dryRun = dryRun
	return &b
}

func newKubectlCommand(args ...string) *kubectlBuilder {
	b := new(kubectlBuilder)
	b.cmd = kubectlCmd(args...)
	return b
}

func (b kubectlBuilder) withStdinData(data string) *kubectlBuilder {
	b.cmd.Stdin = strings.NewReader(data)
	return &b
}

func (b kubectlBuilder) exec() (string, error) {
	if b.dryRun {
		return b.execDryRun()
	}
	var stdout, stderr bytes.Buffer
	cmd := b.cmd
	cmd.Stdout, cmd.Stderr = &stdout, &stderr

	if _, err := exec.LookPath("kubectl"); err != nil {
		return "", fmt.Errorf("Could not find 'kubectl' on $PATH: %s", err)
	}

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func (b kubectlBuilder) execDryRun() (string, error) {
	var stdin string

	buf := new(bytes.Buffer)
	buf.ReadFrom(b.cmd.Stdin)
	stdin = buf.String()

	return fmt.Sprintf("Running '%s %s %s'", b.cmd.Path, strings.Join(b.cmd.Args[1:], " "), stdin), nil // skip arg[0] as it is printed separately
}

// runKubectl is a convenience wrapper over kubectlBuilder
func runKubectl(args ...string) (string, error) {
	return newKubectlCommand(args...).exec()
}

// runKubectlInput is a convenience wrapper over kubectlBuilder that takes input to stdin
func runKubectlInput(data string, args ...string) (string, error) {
	return newKubectlCommand(args...).withStdinData(data).exec()
}
