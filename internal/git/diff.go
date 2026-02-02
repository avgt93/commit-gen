// Package git handles git operations like diff and commit messages.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}
	return string(output), nil
}

func GetRepositoryRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository or failed to get root: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func GetRepositoryName() (string, error) {
	root, err := GetRepositoryRoot()
	if err != nil {
		return "", err
	}
	return filepath.Base(root), nil
}

func GetStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git status: %w", err)
	}
	return string(output), nil
}

func HasStagedChanges() (bool, error) {
	diff, err := GetStagedDiff()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(diff)) > 0, nil
}

func GetChangedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--name-only")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	var result []string
	for _, f := range files {
		if f != "" {
			result = append(result, f)
		}
	}
	return result, nil
}

func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

func GetCommitMessageFile() (string, error) {
	root, err := GetRepositoryRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, ".git", "COMMIT_EDITMSG"), nil
}

func WriteCommitMessage(message string) error {
	msgFile, err := GetCommitMessageFile()
	if err != nil {
		return err
	}

	return os.WriteFile(msgFile, []byte(message), 0o644)
}

func ReadCommitMessage() (string, error) {
	msgFile, err := GetCommitMessageFile()
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(msgFile); os.IsNotExist(err) {
		return "", nil
	}

	content, err := os.ReadFile(msgFile)
	if err != nil {
		return "", fmt.Errorf("failed to read commit message file: %w", err)
	}

	return strings.TrimSpace(string(content)), nil
}
