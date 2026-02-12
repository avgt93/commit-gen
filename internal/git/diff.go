package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const DefaultMaxDiffSize = 32 * 1024

/**
 * DiffResult contains the diff and metadata about whether it was summarized.
 */
type DiffResult struct {
	Diff         string
	IsSummarized bool
	OriginalSize int
}

/**
 * GetStagedDiff returns the staged git diff as a string.
 *
 * @returns The staged diff output
 * @returns An error if the git command fails
 */
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}
	return string(output), nil
}

/**
 * GetStagedDiffStat returns the diff stat showing file change statistics.
 *
 * @returns The diff stat output showing insertions/deletions per file
 * @returns An error if the git command fails
 */
func GetStagedDiffStat() (string, error) {
	cmd := exec.Command("git", "diff", "--staged", "--stat")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff stat: %w", err)
	}
	return string(output), nil
}

/**
 * GetStagedDiffWithLimit returns the staged diff, automatically summarizing
 * if it exceeds the specified maximum size.
 *
 * @param maxSize - Maximum size in bytes before summarizing (0 uses default)
 * @returns A DiffResult containing the diff and metadata about summarization
 * @returns An error if the git command fails
 */
func GetStagedDiffWithLimit(maxSize int) (*DiffResult, error) {
	if maxSize <= 0 {
		maxSize = DefaultMaxDiffSize
	}

	diff, err := GetStagedDiff()
	if err != nil {
		return nil, err
	}

	originalSize := len(diff)

	if originalSize <= maxSize {
		return &DiffResult{
			Diff:         diff,
			IsSummarized: false,
			OriginalSize: originalSize,
		}, nil
	}

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

func summarizeDiff(diff string, maxSize int) (string, error) {
	stat, err := GetStagedDiffStat()
	if err != nil {
		stat = "(unable to get diff stat)"
	}

	files, err := GetChangedFiles()
	if err != nil {
		files = []string{"(unable to get file list)"}
	}

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

	headerSize := sb.Len()
	remainingSpace := maxSize - headerSize - 200

	if remainingSpace > 0 {
		sb.WriteString("=== TRUNCATED DIFF ===\n")
		truncated := truncateDiffSmart(diff, remainingSpace)
		sb.WriteString(truncated)
		sb.WriteString("\n\n... [truncated] ...\n")
	}

	return sb.String(), nil
}

func truncateDiffSmart(diff string, maxLen int) string {
	if len(diff) <= maxLen {
		return diff
	}

	truncated := diff[:maxLen]

	lastHunk := strings.LastIndex(truncated, "\n@@")
	if lastHunk > maxLen/2 {
		hunkEnd := strings.Index(truncated[lastHunk+1:], "\n")
		if hunkEnd > 0 {
			truncated = truncated[:lastHunk+1+hunkEnd]
		}
	} else {
		lastNewline := strings.LastIndex(truncated, "\n")
		if lastNewline > 0 {
			truncated = truncated[:lastNewline]
		}
	}

	return truncated
}

/**
 * GetRepositoryRoot returns the root directory of the current git repository.
 *
 * @returns The absolute path to the repository root
 * @returns An error if not in a git repository
 */
func GetRepositoryRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository or failed to get root: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

/**
 * GetRepositoryName returns the name of the current repository (directory name).
 *
 * @returns The repository name
 * @returns An error if not in a git repository
 */
func GetRepositoryName() (string, error) {
	root, err := GetRepositoryRoot()
	if err != nil {
		return "", err
	}
	return filepath.Base(root), nil
}

/**
 * GetStatus returns the current git status in porcelain format.
 *
 * @returns The git status output
 * @returns An error if the git command fails
 */
func GetStatus() (string, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get git status: %w", err)
	}
	return string(output), nil
}

/**
 * HasStagedChanges checks if there are any staged changes in the repository.
 *
 * @returns true if there are staged changes, false otherwise
 * @returns An error if checking fails
 */
func HasStagedChanges() (bool, error) {
	diff, err := GetStagedDiff()
	if err != nil {
		return false, err
	}
	return len(strings.TrimSpace(diff)) > 0, nil
}

/**
 * GetChangedFiles returns the list of files with staged changes.
 *
 * @returns A slice of file paths with staged changes
 * @returns An error if the git command fails
 */
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

/**
 * IsGitRepository checks if the current directory is inside a git repository.
 *
 * @returns true if in a git repository, false otherwise
 */
func IsGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

/**
 * GetCommitMessageFile returns the path to the git commit message file.
 *
 * @returns The path to .git/COMMIT_EDITMSG
 * @returns An error if not in a git repository
 */
func GetCommitMessageFile() (string, error) {
	root, err := GetRepositoryRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, ".git", "COMMIT_EDITMSG"), nil
}

/**
 * WriteCommitMessage writes a commit message to the git commit message file.
 *
 * @param message - The commit message to write
 * @returns An error if writing fails
 */
func WriteCommitMessage(message string) error {
	msgFile, err := GetCommitMessageFile()
	if err != nil {
		return err
	}

	return os.WriteFile(msgFile, []byte(message), 0o644)
}

/**
 * ChangeEditor sets the git core.editor configuration.
 *
 * @param editor - The editor command to set
 * @returns An error if the git command fails
 */
func ChangeEditor(editor string) error {
	cmd := exec.Command("git", "config", "core.editor", editor)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to change editor: %w", err)
	}
	return nil
}

/**
 * ReadCommitMessage reads the current commit message from the git commit message file.
 *
 * @returns The commit message content, or empty string if file doesn't exist
 * @returns An error if reading fails
 */
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
