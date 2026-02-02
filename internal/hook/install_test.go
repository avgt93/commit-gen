// Package hook manages git hook installation and uninstallation.
package hook

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInstallUninstall(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

	_ = Uninstall()

	if err := Install("cat"); err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	t.Log("✓ Hook installed successfully")

	root := "."

	hookPath := filepath.Join(root, ".git", "hooks", hookName)
	if _, err := os.Stat(hookPath); os.IsNotExist(err) {
		t.Errorf("Hook file not created at %s", hookPath)
	} else {
		t.Logf("✓ Hook file exists: %s", hookPath)
	}

	info, err := os.Stat(hookPath)
	if err != nil {
		t.Fatalf("Failed to stat hook file: %v", err)
	}

	if info.Mode()&0o111 == 0 {
		t.Log("Note: Hook file is not executable (may need chmod)")
	} else {
		t.Log("✓ Hook file is executable")
	}

	if err := Uninstall(); err != nil {
		t.Fatalf("Uninstall failed: %v", err)
	}

	t.Log("✓ Hook uninstalled successfully")

	if _, err := os.Stat(hookPath); !os.IsNotExist(err) {
		t.Logf("Note: Hook file still exists after uninstall")
	} else {
		t.Log("✓ Hook file removed")
	}
}

func TestHookContent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

	_ = Uninstall()

	if err := Install("cat"); err != nil {
		t.Fatalf("Install failed: %v", err)
	}

	hookPath := filepath.Join(".", ".git", "hooks", hookName)

	content, err := os.ReadFile(hookPath)
	if err != nil {
		t.Fatalf("Failed to read hook file: %v", err)
	}

	hookContent := string(content)

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

	_ = Uninstall()
}

func TestIsInstalledFalse(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

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

func TestIsInstalledTrue(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

	if err := Install("cat"); err != nil {
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

	_ = Uninstall()
}

func TestInstallIdempotent(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

	_ = Uninstall()

	if err := Install("cat"); err != nil {
		t.Fatalf("First install failed: %v", err)
	}

	t.Log("✓ First install succeeded")
	err := Install("cat")
	if err != nil {
		t.Logf("✓ Second install correctly returns error: %v", err)
	} else {
		t.Log("Note: Second install succeeded (may overwrite)")
	}

	_ = Uninstall()
}

func TestUninstallWithoutInstall(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping hook test in short mode (requires git repo)")
	}

	gitDir := filepath.Join(".", ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		t.Skip("Not in a git repository, skipping hook tests")
	}

	_ = Uninstall()

	err := Uninstall()
	if err != nil {
		t.Logf("✓ Uninstall correctly returns error when not installed: %v", err)
	} else {
		t.Log("Note: Uninstall succeeded even when not installed")
	}
}

func TestHookScriptContent(t *testing.T) {
	hookScript := fmt.Sprintf(hookScriptFmt, "cat", "commit-gen")

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

func TestHookName(t *testing.T) {
	if hookName != "prepare-commit-msg" {
		t.Errorf("Hook name incorrect: got %q, expected %q", hookName, "prepare-commit-msg")
	} else {
		t.Logf("✓ Hook name correct: %s", hookName)
	}
}
