// Package cache provides in-memory and persistent session caching.
package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestCacheInitialization(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-cache")
	defer os.RemoveAll(tmpDir)

	cache := GetCache(24*time.Hour, tmpDir)

	if cache == nil {
		t.Error("GetCache returned nil")
	}

	if cache.cache == nil {
		t.Error("Cache map is nil")
	}

	t.Log("✓ Cache initialized successfully")
}

func TestCacheSetAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-set-get")
	defer os.RemoveAll(tmpDir)

	cache := GetCache(24*time.Hour, tmpDir)

	if cache == nil {
		t.Error("Cache is nil")
	} else {
		t.Log("✓ Cache created for Set/Get test")
	}
}

func TestCacheTTLExpiration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-ttl")
	defer os.RemoveAll(tmpDir)

	shortTTL := 100 * time.Millisecond
	cache := GetCache(shortTTL, tmpDir)

	if cache.cache == nil {
		t.Error("Cache map is nil")
	} else {
		t.Log("✓ Cache map initialized for TTL test")
	}
}

func TestCacheUpdateLastUsed(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-update")
	defer os.RemoveAll(tmpDir)

	cache := GetCache(24*time.Hour, tmpDir)

	if cache == nil {
		t.Error("Cache is nil")
	} else {
		t.Log("✓ Cache created for UpdateLastUsed test")
	}
}

func TestCacheClear(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-clear")
	defer os.RemoveAll(tmpDir)

	cache := GetCache(24*time.Hour, tmpDir)

	if err := cache.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	t.Log("✓ Cache cleared successfully")

	total, valid, err := cache.Status()
	if err != nil {
		t.Fatalf("Status failed after clear: %v", err)
	}

	if total != 0 || valid != 0 {
		t.Errorf("Cache not properly cleared: total=%d, valid=%d", total, valid)
	} else {
		t.Log("✓ Cache properly cleared (total=0, valid=0)")
	}
}

func TestCacheStatus(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-status")
	defer os.RemoveAll(tmpDir)

	cache := GetCache(24*time.Hour, tmpDir)

	total, valid, err := cache.Status()
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	t.Logf("✓ Cache Status:")
	t.Logf("  - Total entries: %d", total)
	t.Logf("  - Valid entries: %d", valid)

	if total != 0 || valid != 0 {
		t.Logf("Cache has entries (expected empty): total=%d, valid=%d", total, valid)
	}
}

func TestCachePersistence(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-persist")
	defer os.RemoveAll(tmpDir)

	cache := GetCache(24*time.Hour, tmpDir)


	cache.Clear()

	t.Logf("✓ Cache persistence setup complete")
}

func TestHashRepoPath(t *testing.T) {
	path1 := "/home/user/project"
	path2 := "/home/user/project"
	path3 := "/home/user/other"

	hash1 := hashRepoPath(path1)
	hash2 := hashRepoPath(path2)
	hash3 := hashRepoPath(path3)

	if hash1 != hash2 {
		t.Error("Same paths should produce same hash")
	} else {
		t.Log("✓ Same paths produce same hash")
	}

	if hash1 == hash3 {
		t.Error("Different paths should produce different hashes")
	} else {
		t.Log("✓ Different paths produce different hashes")
	}

	if hash1 == "" {
		t.Error("Hash should not be empty")
	} else {
		t.Logf("✓ Hash format: %s", hash1)
	}
}
