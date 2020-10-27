package helpers

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// Run executes spawn process and returns error or nil.
func Run(name string, args []string, wd string) error {
	cmd := exec.Command(name, args...)
	out := bytes.NewBuffer([]byte{})
	cmd.Stdout = out
	cmd.Stderr = out
	if wd != "" {
		cmd.Dir = wd
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%s command failed: %v: %s",
			name, err, strings.TrimRight(out.String(), "\n"))
	}
	return nil
}
