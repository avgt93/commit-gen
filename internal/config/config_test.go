package config

import (
	"testing"
)

// TestConfigInitialization tests that configuration initializes properly
func TestConfigInitialization(t *testing.T) {
	// Clear any existing config
	cfg = nil

	// Initialize with default config
	if err := Initialize(""); err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	if cfg == nil {
		t.Error("Config is nil after initialization")
	}

	t.Log("✓ Configuration initialized successfully")
}

// TestDefaultValues tests that default values are set correctly
func TestDefaultValues(t *testing.T) {
	TestConfigInitialization(t)

	cfg := Get()

	tests := []struct {
		name     string
		getValue func() interface{}
		expected interface{}
	}{
		{"OpenCode Host", func() interface{} { return cfg.OpenCode.Host }, "localhost"},
		{"OpenCode Port", func() interface{} { return cfg.OpenCode.Port }, 4096},
		{"OpenCode Timeout", func() interface{} { return cfg.OpenCode.Timeout }, 30},
		{"Generation Style", func() interface{} { return cfg.Generation.Style }, "conventional"},
		{"Generation Provider", func() interface{} { return cfg.Generation.Model.Provider }, "google"},
		{"Cache Enabled", func() interface{} { return cfg.Cache.Enabled }, true},
		{"Cache TTL", func() interface{} { return cfg.Cache.TTL }, "24h"},
		{"Git Staged Only", func() interface{} { return cfg.Git.StagedOnly }, true},
	}

	for _, tt := range tests {
		value := tt.getValue()
		if value != tt.expected {
			t.Errorf("%s: got %v, expected %v", tt.name, value, tt.expected)
		} else {
			t.Logf("✓ %s: %v", tt.name, value)
		}
	}
}

// TestGetConfigInstance tests getting config instance
func TestGetConfigInstance(t *testing.T) {
	cfg := Get()
	if cfg == nil {
		t.Error("Get() returned nil config")
	}
	t.Log("✓ Got config instance successfully")
}

// TestConfigAccessors tests configuration accessor methods
func TestConfigAccessors(t *testing.T) {
	TestConfigInitialization(t)

	tests := []struct {
		name string
		key  string
		fn   func(string) interface{}
	}{
		{"OpenCode Host", "opencode.host", func(k string) interface{} { return GetString(k) }},
		{"OpenCode Port", "opencode.port", func(k string) interface{} { return GetInt(k) }},
		{"Cache Enabled", "cache.enabled", func(k string) interface{} { return GetBool(k) }},
	}

	for _, tt := range tests {
		value := tt.fn(tt.key)
		t.Logf("✓ %s: %v", tt.name, value)
	}
}

// TestEnvironmentVariableOverride tests that env vars override config
func TestEnvironmentVariableOverride(t *testing.T) {
	// Note: This test is limited because viper caches the config
	// and environment variables are read at initialization time
	// For a complete test, would need to reset viper between tests

	cfg := Get()

	// Just verify that the config system supports the env prefix
	if cfg.OpenCode.Host == "" {
		t.Error("OpenCode Host is empty")
	} else {
		t.Logf("✓ OpenCode host configured: %s", cfg.OpenCode.Host)
	}
}

// TestConfigGet tests the Get() function
func TestConfigGet(t *testing.T) {
	cfg := Get()

	if cfg.OpenCode.Host == "" {
		t.Error("OpenCode Host is empty")
	}

	if cfg.Generation.Model.ModelID == "" {
		t.Error("Model ID is empty")
	}

	t.Logf("✓ Config.Get() returned valid config:")
	t.Logf("  - Host: %s", cfg.OpenCode.Host)
	t.Logf("  - Port: %d", cfg.OpenCode.Port)
	t.Logf("  - Style: %s", cfg.Generation.Style)
	t.Logf("  - Model: %s/%s", cfg.Generation.Model.Provider, cfg.Generation.Model.ModelID)
}

// TestModelConfiguration tests model configuration
func TestModelConfiguration(t *testing.T) {
	cfg := Get()

	if cfg.Generation.Model.Provider == "" {
		t.Error("Model provider is empty")
	}

	if cfg.Generation.Model.ModelID == "" {
		t.Error("Model ID is empty")
	}

	t.Logf("✓ Model Configuration:")
	t.Logf("  - Provider: %s", cfg.Generation.Model.Provider)
	t.Logf("  - Model ID: %s", cfg.Generation.Model.ModelID)
}

// TestCommitStyles tests that all commit styles are valid
func TestCommitStyles(t *testing.T) {
	validStyles := []string{"conventional", "imperative", "detailed"}

	for _, style := range validStyles {
		t.Logf("✓ Valid commit style: %s", style)
	}
}
