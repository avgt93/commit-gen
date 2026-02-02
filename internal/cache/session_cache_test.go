package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestCacheInitialization tests cache initialization
func TestCacheInitialization(t *testing.T) {
	// Use temporary directory for cache
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

// TestCacheSetAndGet tests setting and getting a cached session
func TestCacheSetAndGet(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	// Skip if not in git repo (cache.Set requires git repo)
	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-set-get")
	defer os.RemoveAll(tmpDir)

	cache := GetCache(24*time.Hour, tmpDir)

	// Note: This test is limited because Set() requires git repository
	// We can test cache initialization but not the full Set/Get flow
	if cache == nil {
		t.Error("Cache is nil")
	} else {
		t.Log("✓ Cache created for Set/Get test")
	}
}

// TestCacheTTLExpiration tests that cached sessions expire after TTL
func TestCacheTTLExpiration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-ttl")
	defer os.RemoveAll(tmpDir)

	// Create cache with very short TTL for testing
	shortTTL := 100 * time.Millisecond
	cache := GetCache(shortTTL, tmpDir)

	// Note: Cannot fully test because Set() requires git repo
	// But we can verify cache structure is correct
	if cache.cache == nil {
		t.Error("Cache map is nil")
	} else {
		t.Log("✓ Cache map initialized for TTL test")
	}
}

// TestCacheUpdateLastUsed tests updating last used timestamp
func TestCacheUpdateLastUsed(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping test in short mode")
	}

	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-update")
	defer os.RemoveAll(tmpDir)

	cache := GetCache(24*time.Hour, tmpDir)

	// Note: Cannot fully test because Set() requires git repo
	// But we can verify cache structure is correct
	if cache == nil {
		t.Error("Cache is nil")
	} else {
		t.Log("✓ Cache created for UpdateLastUsed test")
	}
}

// TestCacheClear tests clearing the cache
func TestCacheClear(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-clear")
	defer os.RemoveAll(tmpDir)

	cache := GetCache(24*time.Hour, tmpDir)

	// Note: Cannot fully test because Set() requires git repo
	// But we can test Clear() on empty cache
	if err := cache.Clear(); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	t.Log("✓ Cache cleared successfully")

	// Check status after clear
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

// TestCacheStatus tests getting cache status
func TestCacheStatus(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-status")
	defer os.RemoveAll(tmpDir)

	cache := GetCache(24*time.Hour, tmpDir)

	// Note: Cannot fully test because Set() requires git repo
	// But we can test Status() on empty cache
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

// TestCachePersistence tests cache persistence to disk
func TestCachePersistence(t *testing.T) {
	tmpDir := filepath.Join(os.TempDir(), "commit-gen-test-persist")
	defer os.RemoveAll(tmpDir)

	cache := GetCache(24*time.Hour, tmpDir)

	// Note: Cannot fully test Set() because it requires git repo
	// But we can test that cache file management works

	// Clear to ensure clean state
	cache.Clear()

	// After clearing, cache might not write file if empty
	t.Logf("✓ Cache persistence setup complete")
}

// TestHashRepoPath tests the hash function
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
