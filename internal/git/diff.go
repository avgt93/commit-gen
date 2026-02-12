// Package git handles git operations like diff and commit messages.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DefaultMaxDiffSize is the default maximum diff size in bytes before summarizing.
const DefaultMaxDiffSize = 32 * 1024 // 32KB

// DiffResult contains the diff and metadata about whether it was summarized.
type DiffResult struct {
	Diff         string
	IsSummarized bool
	OriginalSize int
}

func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}
	return string(output), nil
}

// GetStagedDiffStat returns the diff stat (summary of changes).
func GetStagedDiffStat() (string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--stat")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff stat: %w", err)
	}
	return string(output), nil
}

// GetStagedDiffWithLimit returns the diff, summarizing if it exceeds maxSize bytes.
func GetStagedDiffWithLimit(maxSize int) (*DiffResult, error) {
	if maxSize <= 0 {
		maxSize = DefaultMaxDiffSize
	}

	diff, err := GetStagedDiff()
	if err != nil {
		return nil, err
	}

	originalSize := len(diff)

	// If diff is within limit, return as-is
	if originalSize <= maxSize {
		return &DiffResult{
			Diff:         diff,
			IsSummarized: false,
			OriginalSize: originalSize,
		}, nil
	}

	// Diff is too large, create a summary
	summarized, err := summarizeDiff(diff, maxSize)
	if err != nil {
		return nil, err
	}

	return &DiffResult{
		Diff:         summarized,
		IsSummarized: true,
		OriginalSize: originalSize,
	}, nil
}

// summarizeDiff creates a condensed version of a large diff.
func summarizeDiff(diff string, maxSize int) (string, error) {
	// Get the stat summary
	stat, err := GetStagedDiffStat()
	if err != nil {
		stat = "(unable to get diff stat)"
	}

	// Get list of changed files
	files, err := GetChangedFiles()
	if err != nil {
		files = []string{"(unable to get file list)"}
	}

	// Build summary header
	var sb strings.Builder
	sb.WriteString("=== DIFF SUMMARY (original too large) ===\n\n")
	sb.WriteString(fmt.Sprintf("Original diff size: %d bytes\n", len(diff)))
	sb.WriteString(fmt.Sprintf("Files changed: %d\n\n", len(files)))

	sb.WriteString("=== CHANGED FILES ===\n")
	for _, f := range files {
		sb.WriteString(fmt.Sprintf("  - %s\n", f))
	}
	sb.WriteString("\n")

	sb.WriteString("=== DIFF STAT ===\n")
	sb.WriteString(stat)
	sb.WriteString("\n")

	// Calculate how much space we have left for actual diff content
	headerSize := sb.Len()
	remainingSpace := maxSize - headerSize - 200 // Leave buffer for footer

	if remainingSpace > 0 {
		sb.WriteString("=== TRUNCATED DIFF ===\n")
		truncated := truncateDiffSmart(diff, remainingSpace)
		sb.WriteString(truncated)
		sb.WriteString("\n\n... [truncated] ...\n")
	}

	return sb.String(), nil
}

// truncateDiffSmart truncates the diff at a sensible boundary (end of a hunk).
func truncateDiffSmart(diff string, maxLen int) string {
	if len(diff) <= maxLen {
		return diff
	}

	// Try to cut at the end of a hunk (line starting with @@)
	truncated := diff[:maxLen]

	// Find the last complete hunk
	lastHunk := strings.LastIndex(truncated, "\n@@")
	if lastHunk > maxLen/2 {
		// Find the end of this hunk header line
		hunkEnd := strings.Index(truncated[lastHunk+1:], "\n")
		if hunkEnd > 0 {
			truncated = truncated[:lastHunk+1+hunkEnd]
		}
	} else {
		// Just cut at last newline
		lastNewline := strings.LastIndex(truncated, "\n")
		if lastNewline > 0 {
			truncated = truncated[:lastNewline]
		}
	}

	return truncated
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

func ChangeEditor(editor string) error {
	cmd := exec.Command("git", "config", "core.editor", editor)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to change editor: %w", err)
	}
	return nil
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
