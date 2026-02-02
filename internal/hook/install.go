package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/avgt93/commit-gen/internal/git"
)

const hookName = "prepare-commit-msg"

// hookScriptFmt is the content of the git hook (format string)
const hookScriptFmt = `#!/bin/bash
# commit-gen git hook
# Auto-generates commit messages for empty commit messages

MESSAGE_FILE=$1
COMMIT_SOURCE=$2
SHA1=$3

# Only run for normal commits (not for merge commits, etc.)
if [ "$COMMIT_SOURCE" != "" ]; then
  exit 0
fi

# Read the current message and filter out comment lines (starting with #)
MESSAGE=$(grep -v '^#' "$MESSAGE_FILE" 2>/dev/null | xargs)

# Check if message is empty (only whitespace and comments)
if [ -z "$MESSAGE" ]; then
  # Change to git root directory to ensure git commands work
  GIT_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
  if [ -z "$GIT_ROOT" ]; then
    exit 0
  fi
  cd "$GIT_ROOT" || exit 0
  
  # Generate commit message
  TMPFILE=$(mktemp)
  trap "rm -f $TMPFILE" EXIT
  
  if "%s" generate --hook > "$TMPFILE" 2>&1; then
    # Only write if we got output
    if [ -s "$TMPFILE" ]; then
      cat "$TMPFILE" > "$MESSAGE_FILE"
    fi
  fi
fi

exit 0
`

// Install installs the git hook in the current repository
func Install() error {
	root, err := git.GetRepositoryRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	// Get absolute path to the current executable
	exe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	exePath, err := filepath.Abs(exe)
	if err != nil {
		return fmt.Errorf("failed to get absolute executable path: %w", err)
	}

	hookPath := filepath.Join(root, ".git", "hooks", hookName)

	// Create hooks directory if it doesn't exist
	hooksDir := filepath.Dir(hookPath)
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	// Check if hook already exists
	if _, err := os.Stat(hookPath); err == nil {
		// Hook exists, check if it's ours
		content, err := os.ReadFile(hookPath)
		if err == nil && strings.Contains(string(content), "commit-gen") {
			return fmt.Errorf("hook already installed at %s", hookPath)
		}
		return fmt.Errorf("hook already exists at %s (not installed by commit-gen)", hookPath)
	}

	// Format the hook script with the absolute path to the executable
	hookContent := fmt.Sprintf(hookScriptFmt, exePath)

	// Write the hook
	if err := os.WriteFile(hookPath, []byte(hookContent), 0o755); err != nil {
		return fmt.Errorf("failed to write hook: %w", err)
	}

	return nil
}

// Uninstall removes the git hook from the current repository
func Uninstall() error {
	root, err := git.GetRepositoryRoot()
	if err != nil {
		return fmt.Errorf("not in a git repository: %w", err)
	}

	hookPath := filepath.Join(root, ".git", "hooks", hookName)

	// Check if hook exists
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		return fmt.Errorf("hook not found at %s", hookPath)
	}

	// Check if it's our hook
	content, err := os.ReadFile(hookPath)
	if err != nil {
		return fmt.Errorf("failed to read hook: %w", err)
	}

	if !strings.Contains(string(content), "commit-gen") {
		return fmt.Errorf("hook at %s is not a commit-gen hook", hookPath)
	}

	// Remove the hook
	if err := os.Remove(hookPath); err != nil {
		return fmt.Errorf("failed to remove hook: %w", err)
	}

	return nil
}

// IsInstalled checks if the hook is installed
func IsInstalled() (bool, error) {
	root, err := git.GetRepositoryRoot()
	if err != nil {
		return false, err
	}

	hookPath := filepath.Join(root, ".git", "hooks", hookName)

	content, err := os.ReadFile(hookPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	return strings.Contains(string(content), "commit-gen"), nil
}
