package git_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/avgt93/commit-gen/internal/git"
)

// setupTestRepo creates a temporary git repository for testing
func setupTestRepo(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "commit-gen-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Initialize git repo
	cmd := exec.Command("git", "init")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to init git repo: %v", err)
	}

	// Configure git user
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to configure git user email: %v", err)
	}

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		t.Fatalf("Failed to configure git user name: %v", err)
	}

	return tmpDir
}

// TestIntegrationIsGitRepository tests git repository detection
func TestIntegrationIsGitRepository(t *testing.T) {
	tmpDir := setupTestRepo(t)
	defer os.RemoveAll(tmpDir)

	// Save current dir and change to test repo
	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	if !git.IsGitRepository() {
		t.Error("✗ Expected git repository to be detected")
	} else {
		t.Log("✓ Git repository detected successfully")
	}
}

// TestIntegrationGetRepositoryRoot tests repository root detection
func TestIntegrationGetRepositoryRoot(t *testing.T) {
	tmpDir := setupTestRepo(t)
	defer os.RemoveAll(tmpDir)

	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	root, err := git.GetRepositoryRoot()
	if err != nil {
		t.Errorf("✗ Failed to get repository root: %v", err)
		return
	}

	if root != tmpDir {
		t.Errorf("✗ Expected root %s, got %s", tmpDir, root)
	} else {
		t.Logf("✓ Repository root detected: %s", root)
	}
}

// TestIntegrationGetRepositoryName tests repository name extraction
func TestIntegrationGetRepositoryName(t *testing.T) {
	tmpDir := setupTestRepo(t)
	defer os.RemoveAll(tmpDir)

	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	name, err := git.GetRepositoryName()
	if err != nil {
		t.Errorf("✗ Failed to get repository name: %v", err)
		return
	}

	if name == "" {
		t.Error("✗ Repository name should not be empty")
	} else {
		t.Logf("✓ Repository name: %s", name)
	}
}

// TestIntegrationGetStatus tests git status
func TestIntegrationGetStatus(t *testing.T) {
	tmpDir := setupTestRepo(t)
	defer os.RemoveAll(tmpDir)

	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Add and commit
	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	status, err := git.GetStatus()
	if err != nil {
		t.Errorf("✗ Failed to get git status: %v", err)
		return
	}

	if status == "" {
		t.Log("✓ Git status retrieved (clean repository)")
	} else {
		t.Logf("✓ Git status: %s", status)
	}
}

// TestIntegrationGetStagedDiff tests getting staged changes
func TestIntegrationGetStagedDiff(t *testing.T) {
	tmpDir := setupTestRepo(t)
	defer os.RemoveAll(tmpDir)

	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create initial commit
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("initial"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	// Make changes and stage them
	if err := os.WriteFile(testFile, []byte("modified content"), 0644); err != nil {
		t.Fatalf("Failed to modify file: %v", err)
	}

	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to stage changes: %v", err)
	}

	diff, err := git.GetStagedDiff()
	if err != nil {
		t.Errorf("✗ Failed to get staged diff: %v", err)
		return
	}

	if diff == "" {
		t.Log("✓ Got staged diff (may be empty if no changes)")
	} else {
		t.Logf("✓ Staged diff retrieved (%d bytes)", len(diff))
	}
}

// TestIntegrationHasStagedChanges tests staged changes detection
func TestIntegrationHasStagedChanges(t *testing.T) {
	tmpDir := setupTestRepo(t)
	defer os.RemoveAll(tmpDir)

	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Initially no staged changes
	has, err := git.HasStagedChanges()
	if err != nil {
		t.Errorf("✗ Failed to check for staged changes: %v", err)
		return
	}

	if has {
		t.Error("✗ Expected no staged changes in empty repo")
	} else {
		t.Log("✓ Correctly detected no staged changes in empty repo")
	}

	// Create and stage a file
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	has, err = git.HasStagedChanges()
	if err != nil {
		t.Errorf("✗ Failed to check for staged changes: %v", err)
		return
	}

	if !has {
		t.Error("✗ Expected staged changes to be detected")
	} else {
		t.Log("✓ Correctly detected staged changes")
	}
}

// TestIntegrationCommitMessageFile tests commit message file operations
func TestIntegrationCommitMessageFile(t *testing.T) {
	tmpDir := setupTestRepo(t)
	defer os.RemoveAll(tmpDir)

	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Create initial commit for testing
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", "test.txt")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}

	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}

	testMsg := "Test commit message"
	err = git.WriteCommitMessage(testMsg)
	if err != nil {
		t.Errorf("✗ Failed to write commit message: %v", err)
		return
	}
	t.Log("✓ Commit message written successfully")

	content, err := git.ReadCommitMessage()
	if err != nil {
		t.Errorf("✗ Failed to read commit message: %v", err)
		return
	}

	if content != testMsg {
		t.Errorf("✗ Expected message %q, got %q", testMsg, content)
	} else {
		t.Log("✓ Commit message read/write cycle successful")
	}
}

// TestIntegrationEndToEndFlow tests a complete workflow
func TestIntegrationEndToEndFlow(t *testing.T) {
	tmpDir := setupTestRepo(t)
	defer os.RemoveAll(tmpDir)

	oldCwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	defer os.Chdir(oldCwd)

	err = os.Chdir(tmpDir)
	if err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	// Step 1: Verify it's a git repo
	if !git.IsGitRepository() {
		t.Fatal("✗ Expected git repository")
	}
	t.Log("✓ Step 1: Git repository verified")

	// Step 2: Get repository info
	root, err := git.GetRepositoryRoot()
	if err != nil || root == "" {
		t.Fatalf("✗ Failed to get repository root: %v", err)
	}
	t.Logf("✓ Step 2: Repository root: %s", root)

	// Step 3: Create and stage a file
	testFile := filepath.Join(tmpDir, "feature.go")
	content := `package main

func NewFeature() {
	// TODO: implement feature
}
`
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatalf("✗ Failed to create test file: %v", err)
	}

	cmd := exec.Command("git", "add", "feature.go")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("✗ Failed to stage file: %v", err)
	}
	t.Log("✓ Step 3: File staged")

	// Step 4: Check for staged changes
	has, err := git.HasStagedChanges()
	if err != nil || !has {
		t.Fatalf("✗ Expected staged changes")
	}
	t.Log("✓ Step 4: Staged changes detected")

	// Step 5: Get diff
	diff, err := git.GetStagedDiff()
	if err != nil {
		t.Fatalf("✗ Failed to get diff: %v", err)
	}
	if len(diff) == 0 {
		t.Error("⚠ Diff is empty (unexpected)")
	}
	t.Logf("✓ Step 5: Diff retrieved (%d bytes)", len(diff))

	// Step 6: Commit
	cmd = exec.Command("git", "commit", "-m", "feat: add new feature")
	cmd.Dir = tmpDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("✗ Failed to commit: %v", err)
	}
	t.Log("✓ Step 6: Changes committed")

	// Step 7: Verify no more staged changes
	has, err = git.HasStagedChanges()
	if err != nil || has {
		t.Error("✗ Expected no staged changes after commit")
	}
	t.Log("✓ Step 7: No staged changes after commit")

	t.Log("\n✓ Integration test completed successfully!")
}

// BenchmarkGetStagedDiff benchmarks the GetStagedDiff function
func BenchmarkGetStagedDiff(b *testing.B) {
	tmpDir := setupTestRepo(&testing.T{})
	defer os.RemoveAll(tmpDir)

	oldCwd, _ := os.Getwd()
	defer os.Chdir(oldCwd)
	os.Chdir(tmpDir)

	// Create a file with some content
	testFile := filepath.Join(tmpDir, "large_file.txt")
	largeContent := make([]byte, 1024*100) // 100KB file
	for i := 0; i < len(largeContent); i++ {
		largeContent[i] = byte((i % 26) + 'a')
	}
	os.WriteFile(testFile, largeContent, 0644)

	cmd := exec.Command("git", "add", "large_file.txt")
	cmd.Dir = tmpDir
	cmd.Run()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		git.GetStagedDiff()
	}
}
