package kubectl

import (
	"bytes"
	"fmt"
	"os/exec"
)

type Config struct {
	Bin    string
	Server string
	Token  string
}

// Cmd ...
type Cmd struct {
	config *Config
}

func New(c *Config) *Cmd {
	if c.Bin == "" {
		c.Bin = "kubectl"
	}

	return &Cmd{config: c}
}

// Validate ...
func (c *Cmd) Validate(data []byte, arg ...string) ([]byte, error) {
	return c.exec(data, append([]string{"apply", "--dry-run=client", "-f", "-"}, arg...))
}

// Get ...
func (c *Cmd) Get(arg ...string) ([]byte, error) {
	return c.exec(nil, append([]string{"get"}, arg...))
}

// Apply ...
func (c *Cmd) Apply(data []byte, arg ...string) ([]byte, error) {
	return c.exec(data, append([]string{"apply", "-f", "-"}, arg...))
}

// Delete ...
func (c *Cmd) Delete(data []byte, arg ...string) ([]byte, error) {
	return c.exec(data, append([]string{"delete", "-f", "-"}, arg...))
}

// Delete ...
func (c *Cmd) DeleteBySelector(selector string) ([]byte, error) {
	selector = "--selector=" + selector
	return c.Exec(append([]string{"delete", "ns"}, selector))
}

func (c *Cmd) Exec(args []string) ([]byte, error) {
	args = append(args, "--insecure-skip-tls-verify")
	if c.config.Server != "" {
		args = append(args, "--server", c.config.Server)
	}
	if c.config.Token != "" {
		args = append(args, "--token", c.config.Token)
	}
	cmd := exec.Command(c.config.Bin, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("kubectl exec error: %w : stderr: %s", err, output)
	}

	return output, nil
}

func (c *Cmd) exec(data []byte, args []string) ([]byte, error) {
	args = append(args, "--insecure-skip-tls-verify")
	if c.config.Server != "" {
		args = append(args, "--server", c.config.Server)
	}
	if c.config.Token != "" {
		args = append(args, "--token", c.config.Token)
	}

	cmd := exec.Command(c.config.Bin, args...)
	cmd.Stdin = bytes.NewBuffer(data)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("kubectl exec error: %w : stderr: %s", err, output)
	}

	return output, nil
}
