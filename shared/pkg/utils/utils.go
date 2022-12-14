package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

func RunCommand(cmd *exec.Cmd) ([]byte, error) {
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run command %v, error: %w, stdout: %v, stderr: %v", cmd.String(), err, outb.String(), errb.String())
	}
	return outb.Bytes(), nil
}
