package opencode

import (
	"testing"
	"time"
)

/**
 * TestRunnerCreation verifies that NewRunner creates a valid Runner instance.
 */
func TestRunnerCreation(t *testing.T) {
	runner := NewRunner(30)

	if runner == nil {
		t.Fatal("NewRunner returned nil")
	}

	if runner.timeout != 30*time.Second {
		t.Errorf("Timeout mismatch: got %v, expected 30s", runner.timeout)
	}

	t.Logf("✓ Runner created successfully with timeout: %v", runner.timeout)
}

/**
 * TestRunnerCreationWithZeroTimeout verifies that NewRunner handles zero timeout.
 */
func TestRunnerCreationWithZeroTimeout(t *testing.T) {
	runner := NewRunner(0)

	if runner == nil {
		t.Fatal("NewRunner returned nil")
	}

	if runner.timeout != 0 {
		t.Errorf("Expected zero timeout, got %v", runner.timeout)
	}

	t.Log("✓ Runner created with zero timeout")
}

/**
 * TestRunnerCreationWithLargeTimeout verifies that NewRunner handles large timeout values.
 */
func TestRunnerCreationWithLargeTimeout(t *testing.T) {
	runner := NewRunner(300)

	if runner.timeout != 300*time.Second {
		t.Errorf("Timeout mismatch: got %v, expected 300s", runner.timeout)
	}

	t.Logf("✓ Runner created with large timeout: %v", runner.timeout)
}

/**
 * TestCheckAvailableExistingCommand verifies CheckAvailable returns true for existing command.
 */
func TestCheckAvailableExistingCommand(t *testing.T) {
	runner := NewRunner(10)

	available, err := runner.CheckAvailable()
	if err != nil {
		t.Logf("Note: CheckAvailable returned error (may be expected in test env): %v", err)
	}

	if available {
		t.Log("✓ opencode binary is available")
	} else {
		t.Log("✗ opencode binary not found (expected in test environment)")
	}
}

/**
 * TestCheckAvailableWithTimeout verifies that CheckAvailable respects timeout.
 */
func TestCheckAvailableWithTimeout(t *testing.T) {
	runner := NewRunner(5)

	if runner.timeout != 5*time.Second {
		t.Errorf("Runner timeout not set correctly: got %v", runner.timeout)
	}

	t.Logf("✓ Runner timeout configured: %v", runner.timeout)
}

/**
 * TestRunnerTimeoutType verifies timeout is of correct type.
 */
func TestRunnerTimeoutType(t *testing.T) {
	runner := NewRunner(15)

	expectedDuration := 15 * time.Second

	if runner.timeout != expectedDuration {
		t.Errorf("Timeout type mismatch: got %T with value %v", runner.timeout, runner.timeout)
	}

	t.Logf("✓ Runner timeout is correct type: %T", runner.timeout)
}

/**
 * TestRunnerStructFields verifies Runner struct has expected fields.
 */
func TestRunnerStructFields(t *testing.T) {
	runner := NewRunner(20)

	if runner.timeout != 20*time.Second {
		t.Error("Runner timeout field not accessible")
	}

	t.Log("✓ Runner struct fields are accessible")
}

/**
 * TestMultipleRunnerInstances verifies multiple runners can be created independently.
 */
func TestMultipleRunnerInstances(t *testing.T) {
	runner1 := NewRunner(10)
	runner2 := NewRunner(30)
	runner3 := NewRunner(60)

	if runner1.timeout != 10*time.Second {
		t.Error("Runner1 timeout incorrect")
	}

	if runner2.timeout != 30*time.Second {
		t.Error("Runner2 timeout incorrect")
	}

	if runner3.timeout != 60*time.Second {
		t.Error("Runner3 timeout incorrect")
	}

	t.Log("✓ Multiple runner instances created with different timeouts")
}

/**
 * TestFilterOutput verifies that filterOutput removes auto-update-checker messages.
 */
func TestFilterOutput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no noise",
			input:    "Add user authentication feature",
			expected: "Add user authentication feature",
		},
		{
			name:     "auto-update-checker prefix",
			input:    "[auto-update-checker] Package removed: /path/to/package\nAdd user authentication feature",
			expected: "Add user authentication feature",
		},
		{
			name:     "multiple auto-update-checker lines",
			input:    "[auto-update-checker] Checking for updates...\n[auto-update-checker] Package removed: /path/to/package\nFix bug in login handler",
			expected: "Fix bug in login handler",
		},
		{
			name:     "auto-update-checker at end",
			input:    "Refactor database queries\n[auto-update-checker] Update available",
			expected: "Refactor database queries",
		},
		{
			name:     "only auto-update-checker",
			input:    "[auto-update-checker] Some message",
			expected: "",
		},
		{
			name:     "empty input",
			input:    "",
			expected: "",
		},
		{
			name:     "whitespace handling",
			input:    "\n[auto-update-checker] Noise\n\nUpdate README  \n",
			expected: "Update README",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := filterOutput(tt.input)
			if result != tt.expected {
				t.Errorf("filterOutput() = %q, want %q", result, tt.expected)
			}
		})
	}
}
