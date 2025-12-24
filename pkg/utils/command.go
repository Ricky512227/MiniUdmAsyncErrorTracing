package utils

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// ExecuteCommand executes a shell command and returns output
func ExecuteCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("command '%s %s' failed: %w", command, strings.Join(args, " "), err)
	}
	return string(output), nil
}

// ExecuteCommandWithTimeout executes a shell command with a timeout
func ExecuteCommandWithTimeout(timeout time.Duration, command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	
	// Create a channel to signal completion
	done := make(chan error, 1)
	var output []byte
	var err error
	
	go func() {
		output, err = cmd.CombinedOutput()
		done <- err
	}()

	// Wait for completion or timeout
	select {
	case <-time.After(timeout):
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		return "", fmt.Errorf("command timeout after %v", timeout)
	case err := <-done:
		if err != nil {
			return string(output), fmt.Errorf("command '%s %s' failed: %w", command, strings.Join(args, " "), err)
		}
		return string(output), nil
	}
}

