package cache

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/avgt93/commit-gen/internal/git"
)

type CachedSession struct {
	SessionID  string    `json:"session_id"`
	RepoPath   string    `json:"repo_path"`
	CreatedAt  time.Time `json:"created_at"`
	LastUsedAt time.Time `json:"last_used_at"`
}

type SessionCache struct {
	mu       sync.RWMutex
	cache    map[string]*CachedSession
	ttl      time.Duration
	cachedir string
}

var instance *SessionCache
var once sync.Once

// GetCache returns the singleton cache instance
func GetCache(ttl time.Duration, cachedir string) *SessionCache {
	once.Do(func() {
		instance = &SessionCache{
			cache:    make(map[string]*CachedSession),
			ttl:      ttl,
			cachedir: cachedir,
		}
		instance.load()
	})
	return instance
}

// Get retrieves a cached session for the current repository
func (sc *SessionCache) Get() (*CachedSession, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	repoPath, err := git.GetRepositoryRoot()
	if err != nil {
		return nil, err
	}

	key := hashRepoPath(repoPath)
	session, exists := sc.cache[key]
	if !exists {
		return nil, nil
	}

	// Check if session has expired
	if time.Since(session.CreatedAt) > sc.ttl {
		return nil, nil
	}

	return session, nil
}

// Set stores a session in the cache
func (sc *SessionCache) Set(sessionID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	repoPath, err := git.GetRepositoryRoot()
	if err != nil {
		return err
	}

	key := hashRepoPath(repoPath)
	now := time.Now()

	sc.cache[key] = &CachedSession{
		SessionID:  sessionID,
		RepoPath:   repoPath,
		CreatedAt:  now,
		LastUsedAt: now,
	}

	return sc.save()
}

// UpdateLastUsed updates the last used timestamp
func (sc *SessionCache) UpdateLastUsed(sessionID string) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	for _, session := range sc.cache {
		if session.SessionID == sessionID {
			session.LastUsedAt = time.Now()
			return sc.save()
		}
	}

	return fmt.Errorf("session not found in cache")
}

// Clear removes all cached sessions
func (sc *SessionCache) Clear() error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	sc.cache = make(map[string]*CachedSession)
	return sc.save()
}

// Status returns cache statistics
func (sc *SessionCache) Status() (int, int, error) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	totalEntries := len(sc.cache)
	validEntries := 0

	for _, session := range sc.cache {
		if time.Since(session.CreatedAt) <= sc.ttl {
			validEntries++
		}
	}

	return totalEntries, validEntries, nil
}

// Helper functions

func hashRepoPath(path string) string {
	hash := md5.Sum([]byte(path))
	return fmt.Sprintf("%x", hash)
}

func (sc *SessionCache) load() error {
	cacheFile := filepath.Join(sc.cachedir, "sessions.json")

	data, err := os.ReadFile(cacheFile)
	if err != nil {
		// Cache file doesn't exist yet, which is fine
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var cached map[string]*CachedSession
	if err := json.Unmarshal(data, &cached); err != nil {
		return err
	}

	sc.cache = cached
	return nil
}

func (sc *SessionCache) save() error {
	// Ensure cache directory exists
	if err := os.MkdirAll(sc.cachedir, 0o755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	cacheFile := filepath.Join(sc.cachedir, "sessions.json")
	data, err := json.MarshalIndent(sc.cache, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, data, 0o644)
}
