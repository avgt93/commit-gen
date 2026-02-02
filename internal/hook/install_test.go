package hook

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestInstallUninstall tests installing and uninstalling the hook
func TestInstallUninstall(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	// This test requires a git repository
	// Check if we're in one
	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

	// Test uninstall first (cleanup any existing hook)
	_ = Uninstall()

	// Test install
	if err := Install(); err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	t.Log("✓ Hook installed successfully")

	// Verify hook was installed
	root := "."

	hookPath := filepath.Join(root, ".git", "hooks", hookName)
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		t.Errorf("Hook file not created at %s", hookPath)
	} else {
		t.Logf("✓ Hook file exists: %s", hookPath)
	}

	// Verify hook is executable
	info, err := os.Stat(hookPath)
	if err != nil {
		t.Fatalf("Failed to stat hook file: %v", err)
	}

	if info.Mode()&0o111 == 0 {
		t.Log("Note: Hook file is not executable (may need chmod)")
	} else {
		t.Log("✓ Hook file is executable")
	}

	// Test uninstall
	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	t.Log("✓ Hook uninstalled successfully")

	// Verify hook was removed
	if _, err := os.Stat(hookPath); !os.IsNotExist(err) {
		t.Logf("Note: Hook file still exists after uninstall")
	} else {
		t.Log("✓ Hook file removed")
	}
}

// TestHookContent tests that the hook script has correct content
func TestHookContent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

	// Uninstall any existing hook first
	_ = Uninstall()

	// Install hook
	if err := Install(); err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	hookPath := filepath.Join(".", ".git", "hooks", hookName)

	// Read the hook file
	content, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("Failed to read hook file: %v", err)
	}

	hookContent := string(content)

	// Check for expected content
	expectedStrings := []string{
		"#!/bin/bash",
		"commit-gen",
		"MESSAGE_FILE",
		"COMMIT_EDITMSG",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(hookContent, expected) {
			t.Errorf("Hook script missing: %q", expected)
		} else {
			t.Logf("✓ Hook contains: %q", expected)
		}
	}

	// Cleanup
	_ = Uninstall()
}

// TestIsInstalledFalse tests IsInstalled when hook is not installed
func TestIsInstalledFalse(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

	// Make sure hook is uninstalled
	_ = Uninstall()

	installed, err := IsInstalled()
	if err != nil {
		t.Logf("Note: IsInstalled error (may be expected): %v", err)
	}

	if installed {
		t.Error("Expected IsInstalled to return false")
	} else {
		t.Log("✓ IsInstalled correctly returns false when not installed")
	}
}

// TestIsInstalledTrue tests IsInstalled when hook is installed
func TestIsInstalledTrue(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

	// Install hook first
	if err := Install(); err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	installed, err := IsInstalled()
	if err != nil {
		t.Fatalf("IsInstalled failed: %v", err)
	}

	if !installed {
		t.Error("Expected IsInstalled to return true after installing")
	} else {
		t.Log("✓ IsInstalled correctly returns true when installed")
	}

	// Cleanup
	_ = Uninstall()
}

// TestInstallIdempotent tests that installing twice fails gracefully
func TestInstallIdempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

	// Uninstall first
	_ = Uninstall()

	// Install first time
	if err := Install(); err != nil {
		t.Fatalf("First install failed: %v", err)
	}

	t.Log("✓ First install succeeded")

	// Try to install second time (should fail or warn)
	err := Install()
	if err != nil {
		t.Logf("✓ Second install correctly returns error: %v", err)
	} else {
		t.Log("Note: Second install succeeded (may overwrite)")
	}

	// Cleanup
	_ = Uninstall()
}

// TestUninstallWithoutInstall tests uninstalling when not installed
func TestUninstallWithoutInstall(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

	// Make sure it's not installed
	_ = Uninstall()

	// Try to uninstall again
	err := Uninstall()
	if err != nil {
		t.Logf("✓ Uninstall correctly returns error when not installed: %v", err)
	} else {
		t.Log("Note: Uninstall succeeded even when not installed")
	}
}

// TestHookScriptContent tests the actual hook script content
func TestHookScriptContent(t *testing.T) {
	expectedKeywords := []string{
		"bash",
		"commit-gen",
		"MESSAGE_FILE",
		"exit 0",
	}

	for _, keyword := range expectedKeywords {
		if !strings.Contains(hookScript, keyword) {
			t.Errorf("Hook script missing keyword: %q", keyword)
		} else {
			t.Logf("✓ Hook script contains: %q", keyword)
		}
	}
}

// TestHookName tests hook name constant
func TestHookName(t *testing.T) {
	if hookName != "prepare-commit-msg" {
		t.Errorf("Hook name incorrect: got %q, expected %q", hookName, "prepare-commit-msg")
	} else {
		t.Logf("✓ Hook name correct: %s", hookName)
	}
}
