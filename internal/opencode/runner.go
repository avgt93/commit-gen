// Package opencode provides HTTP client and CLI runner for OpenCode.
package opencode

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Runner executes opencode CLI commands directly via subprocess.
type Runner struct {
	timeout time.Duration
}

// NewRunner creates a new Runner with the specified timeout in seconds.
func NewRunner(timeout int) *Runner {
	return &Runner{
		timeout: time.Duration(timeout) * time.Second,
	}
}

// CheckAvailable verifies that the opencode binary is available in PATH.
func (r *Runner) CheckAvailable() (bool, error) {
	_, err := exec.LookPath("opencode")
	if err != nil {
		return false, fmt.Errorf("opencode binary not found in PATH: %w", err)
	}
	return true, nil
}

// Generate runs opencode with the given prompt and returns the generated text.
func (r *Runner) Generate(prompt string, model *Model) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	args := []string{"run"}

	// Model format for opencode run is: --model provider/model_id
	if model != nil && model.ProviderID != "" && model.ModelID != "" {
		args = append(args, "--model", fmt.Sprintf("%s/%s", model.ProviderID, model.ModelID))
	}

	// Add prompt as the message
	args = append(args, prompt)

	cmd := exec.CommandContext(ctx, "opencode", args...)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("opencode run timed out after %v", r.timeout)
		}
		return "", fmt.Errorf("opencode run failed: %w - %s", err, stderr.String())
	}

	return strings.TrimSpace(stdout.String()), nil
}
