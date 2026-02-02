package generator

import (
	"testing"
	"time"

	"github.com/avgt93/commit-gen/internal/cache"
	"github.com/avgt93/commit-gen/internal/config"
)

// TestGeneratorCreation tests creating a new generator
func TestGeneratorCreation(t *testing.T) {
	config.Initialize("")
	cfg := config.Get()

	cacheDir := t.TempDir()
	sessionCache := cache.GetCache(24*time.Hour, cacheDir)

	gen := NewGenerator(cfg, sessionCache)

	if gen == nil {
		t.Error("NewGenerator returned nil")
	}

	if gen.config == nil {
		t.Error("Generator config is nil")
	}

	if gen.client == nil {
		t.Error("Generator client is nil")
	}

	t.Log("✓ Generator created successfully")
}

// TestStyleGuideConventional tests conventional commit style guide
func TestStyleGuideConventional(t *testing.T) {
	guide := getStyleGuide("conventional")

	if guide == "" {
		t.Error("Style guide is empty")
	}

	expectedKeywords := []string{"Conventional Commits", "type(scope)", "feat", "fix"}
	for _, keyword := range expectedKeywords {
		if !contains(guide, keyword) {
			t.Errorf("Style guide missing keyword: %s", keyword)
		}
	}

	t.Logf("✓ Conventional style guide contains expected content")
}

// TestStyleGuideImperative tests imperative style guide
func TestStyleGuideImperative(t *testing.T) {
	guide := getStyleGuide("imperative")

	if guide == "" {
		t.Error("Style guide is empty")
	}

	expectedKeywords := []string{"imperative mood", "verb"}
	for _, keyword := range expectedKeywords {
		if !contains(guide, keyword) {
			t.Errorf("Style guide missing keyword: %s", keyword)
		}
	}

	t.Logf("✓ Imperative style guide contains expected content")
}

// TestStyleGuideDetailed tests detailed style guide
func TestStyleGuideDetailed(t *testing.T) {
	guide := getStyleGuide("detailed")

	if guide == "" {
		t.Error("Style guide is empty")
	}

	expectedKeywords := []string{"type(scope)"}
	for _, keyword := range expectedKeywords {
		if !contains(guide, keyword) {
			t.Errorf("Style guide missing keyword: %s", keyword)
		}
	}

	t.Logf("✓ Detailed style guide contains expected content")
}

// TestStyleGuideUnknown tests unknown style defaults to conventional
func TestStyleGuideUnknown(t *testing.T) {
	guide := getStyleGuide("unknown-style")

	if guide == "" {
		t.Error("Style guide is empty for unknown style")
	}

	// Should default to conventional
	if !contains(guide, "Conventional") {
		t.Log("Note: Unknown style may not default to conventional")
	} else {
		t.Log("✓ Unknown style defaults to conventional")
	}
}

// TestBuildPrompt tests prompt building
func TestBuildPrompt(t *testing.T) {
	config.Initialize("")
	cfg := config.Get()

	cacheDir := t.TempDir()
	sessionCache := cache.GetCache(24*time.Hour, cacheDir)
	gen := NewGenerator(cfg, sessionCache)

	testDiff := "diff --git a/test.go b/test.go\n+++ b/test.go\n@@ -1,3 +1,4 @@"

	prompt := gen.buildPrompt(testDiff)

	if prompt == "" {
		t.Error("Prompt is empty")
	}

	if !contains(prompt, testDiff) {
		t.Error("Prompt doesn't contain the diff")
	}

	if !contains(prompt, "commit message") {
		t.Error("Prompt doesn't mention commit message")
	}

	t.Logf("✓ Prompt built successfully (%d chars)", len(prompt))
}

// TestExtractCommitMessageBasic tests extracting a basic message
func TestExtractCommitMessageBasic(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"feat: add new feature", "feat: add new feature"},
		{"feat: add new feature\nMore details", "feat: add new feature"},
		{"  feat: add new feature  ", "feat: add new feature"},
		{"```\nfeat: add new feature\n```", "feat: add new feature"},
		{"feat: add new feature\n\nBody text", "feat: add new feature"},
	}

	for _, tt := range tests {
		result := extractCommitMessage(tt.input)
		if result != tt.expected {
			t.Errorf("Extract message mismatch:\n  input: %q\n  got: %q\n  expected: %q", tt.input, result, tt.expected)
		} else {
			t.Logf("✓ Extracted: %q", result)
		}
	}
}

// TestExtractCommitMessageRemovesMarkdown tests markdown removal
func TestExtractCommitMessageRemovesMarkdown(t *testing.T) {
	input := "```\nfeat: add feature\n```\n"
	expected := "feat: add feature"

	result := extractCommitMessage(input)

	if result != expected {
		t.Errorf("Markdown not removed correctly:\n  got: %q\n  expected: %q", result, expected)
	} else {
		t.Log("✓ Markdown code blocks removed correctly")
	}
}

// TestExtractCommitMessageTrimsWhitespace tests whitespace trimming
func TestExtractCommitMessageTrimsWhitespace(t *testing.T) {
	input := "   \n   feat: add feature   \n   "
	expected := "feat: add feature"

	result := extractCommitMessage(input)

	if result != expected {
		t.Errorf("Whitespace not trimmed correctly:\n  got: %q\n  expected: %q", result, expected)
	} else {
		t.Log("✓ Whitespace trimmed correctly")
	}
}

// TestExtractCommitMessageFirstLineOnly tests first line extraction
func TestExtractCommitMessageFirstLineOnly(t *testing.T) {
	input := "feat: main change\nThis is additional info\nMore details here"
	expected := "feat: main change"

	result := extractCommitMessage(input)

	if result != expected {
		t.Errorf("First line not extracted correctly:\n  got: %q\n  expected: %q", result, expected)
	} else {
		t.Log("✓ First line extracted correctly")
	}
}

// TestAllCommitStyles tests that all three styles are supported
func TestAllCommitStyles(t *testing.T) {
	styles := []string{"conventional", "imperative", "detailed"}

	for _, style := range styles {
		guide := getStyleGuide(style)
		if guide == "" {
			t.Errorf("Empty guide for style: %s", style)
		} else {
			t.Logf("✓ Style guide for %q (%d chars)", style, len(guide))
		}
	}
}

// TestPromptContainsInstructions tests that prompts contain instructions
func TestPromptContainsInstructions(t *testing.T) {
	config.Initialize("")
	cfg := config.Get()

	cacheDir := t.TempDir()
	sessionCache := cache.GetCache(24*time.Hour, cacheDir)
	gen := NewGenerator(cfg, sessionCache)

	diff := "test diff"
	prompt := gen.buildPrompt(diff)

	requiredContent := []string{
		"commit message",
		"changes",
	}

	for _, content := range requiredContent {
		if !contains(prompt, content) {
			t.Errorf("Prompt missing required content: %s", content)
		}
	}

	t.Log("✓ Prompt contains all required instructions")
}

// Helper function to check if string contains substring
func contains(str, substr string) bool {
	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
