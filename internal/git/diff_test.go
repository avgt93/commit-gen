// Package git handles git operations like diff and commit messages.
package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestIsGitRepository(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}
}

func TestGetRepositoryRoot(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	root, err := GetRepositoryRoot()
	if err != nil {
		t.Fatalf("GetRepositoryRoot failed: %v", err)
	}

	if root == "" {
		t.Error("GetRepositoryRoot returned empty string")
	}

	if _, err := os.Stat(root); os.IsNotExist(err) {
		t.Errorf("Repository root path does not exist: %s", root)
	}

	gitDir := filepath.Join(root, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Errorf(".git directory not found at: %s", gitDir)
	}
}

func TestGetRepositoryName(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	name, err := GetRepositoryName()
	if err != nil {
		t.Fatalf("GetRepositoryName failed: %v", err)
	}

	if name == "" {
		t.Error("GetRepositoryName returned empty string")
	}

	expectedName := "commit-gen"
	if name != expectedName {
		t.Logf("Repository name: %s (expected something like: %s)", name, expectedName)
	}
}

func TestGetStatus(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	status, err := GetStatus()
	if err != nil {
		t.Fatalf("GetStatus failed: %v", err)
	}

	if status == "" {
		t.Log("No uncommitted changes (status is empty)")
	} else {
		t.Logf("Git status: %s", status)
	}
}

func TestGetStagedDiff(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	diff, err := GetStagedDiff()
	if err != nil {
		t.Fatalf("GetStagedDiff failed: %v", err)
	}

	if diff == "" {
		t.Log("No staged changes (diff is empty)")
	} else {
		t.Logf("Staged diff length: %d bytes", len(diff))
		if len(diff) > 100 {
			t.Logf("Diff preview: %s...", diff[:100])
		}
	}
}

func TestGetChangedFiles(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	files, err := GetChangedFiles()
	if err != nil {
		t.Fatalf("GetChangedFiles failed: %v", err)
	}

	t.Logf("Number of staged files: %d", len(files))
	for _, f := range files {
		t.Logf("  - %s", f)
	}
}

func TestHasStagedChanges(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	has, err := HasStagedChanges()
	if err != nil {
		t.Fatalf("HasStagedChanges failed: %v", err)
	}

	if has {
		t.Log("Staged changes detected")
	} else {
		t.Log("No staged changes")
	}
}

func TestCommitMessageFileOperations(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	msgFile, err := GetCommitMessageFile()
	if err != nil {
		t.Fatalf("GetCommitMessageFile failed: %v", err)
	}

	if msgFile == "" {
		t.Error("GetCommitMessageFile returned empty path")
	}

	t.Logf("Commit message file path: %s", msgFile)

	msg, err := ReadCommitMessage()
	if err != nil {
		t.Logf("Note: ReadCommitMessage error (may not exist yet): %v", err)
	} else {
		t.Logf("Current commit message: %q", msg)
	}
}

func TestGitCommandExecution(t *testing.T) {
	if !IsGitRepository() {
		t.Skip("Not in a git repository, skipping test")
	}

	cmd := exec.Command("git", "rev-parse", "--git-dir")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Git command execution failed: %v", err)
	}

	if len(output) == 0 {
		t.Error("Git command returned empty output")
	}

	t.Logf("Git directory: %s", string(output))
}
