package generator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/avgt93/commit-gen/internal/cache"
	"github.com/avgt93/commit-gen/internal/config"
	"github.com/avgt93/commit-gen/internal/git"
	"github.com/avgt93/commit-gen/internal/opencode"
)

var ErrServerNotRunning = errors.New("opencode server is not running")

/**
 * Generator handles commit message generation using either server or run mode.
 */
type Generator struct {
	client *opencode.Client
	runner *opencode.Runner
	cache  *cache.SessionCache
	config *config.Config
	mode   string
}

/**
 * NewGenerator creates a new Generator based on the configured mode.
 *
 * @param cfg - The application configuration
 * @param cacheInstance - The session cache for server mode
 * @returns A new Generator instance
 */
func NewGenerator(cfg *config.Config, cacheInstance *cache.SessionCache) *Generator {
	mode := cfg.OpenCode.Mode
	if mode == "" {
		mode = "run"
	}

	gen := &Generator{
		cache:  cacheInstance,
		config: cfg,
		mode:   mode,
	}

	if mode == "server" {
		gen.client = opencode.NewClient(cfg.OpenCode.Host, cfg.OpenCode.Port, cfg.OpenCode.Timeout)
	} else {
		gen.runner = opencode.NewRunner(cfg.OpenCode.Timeout)
	}

	return gen
}

/**
 * GetMode returns the current operation mode.
 *
 * @returns "run" or "server"
 */
func (g *Generator) GetMode() string {
	return g.mode
}

/**
 * GetConfig returns the generator's configuration.
 *
 * @returns The Config instance
 */
func (g *Generator) GetConfig() *config.Config {
	return g.config
}

/**
 * Generate creates a commit message from staged changes.
 *
 * @returns The generated commit message
 * @returns An error if generation fails
 */
func (g *Generator) Generate() (string, error) {
	maxSize := g.config.Git.MaxDiffSize
	if maxSize <= 0 {
		maxSize = git.DefaultMaxDiffSize
	}

	diffResult, err := git.GetStagedDiffWithLimit(maxSize)
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}

	if strings.TrimSpace(diffResult.Diff) == "" {
		return "", fmt.Errorf("no staged changes found")
	}

	// if diffResult.IsSummarized {
	// return "", fmt.Errorf("note: Large diff (%d bytes) was summarized for AI processing", diffResult.OriginalSize)
	// }

	if g.mode == "server" {
		return g.generateWithServer(diffResult.Diff, diffResult.IsSummarized)
	}
	return g.generateWithRunner(diffResult.Diff, diffResult.IsSummarized)
}

func (g *Generator) generateWithRunner(diff string, isSummarized bool) (string, error) {
	prompt := g.buildPrompt(diff, isSummarized)

	model := &opencode.Model{
		ProviderID: g.config.Generation.Model.Provider,
		ModelID:    g.config.Generation.Model.ModelID,
	}

	response, err := g.runner.Generate(prompt, model)
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	message := extractCommitMessage(response)
	return message, nil
}

func (g *Generator) generateWithServer(diff string, isSummarized bool) (string, error) {
	healthy, err := g.client.CheckHealth()
	if err != nil || !healthy {
		fmt.Printf("%v at %s:%d", ErrServerNotRunning, g.config.OpenCode.Host, g.config.OpenCode.Port)
		return "", fmt.Errorf("failed to start opencode server: %w", err)
	}

	var sessionID string
	cachedSession, err := g.cache.Get()
	if err == nil && cachedSession != nil {
		sessionID = cachedSession.SessionID
	} else {
		repoName, err := git.GetRepositoryName()
		if err != nil {
			repoName = "project"
		}

		session, err := g.client.CreateSession(fmt.Sprintf("commit-gen: %s", repoName))
		if err != nil {
			return "", fmt.Errorf("failed to create OpenCode session: %w", err)
		}

		sessionID = session.ID
		if err := g.cache.Set(sessionID); err != nil {
			fmt.Printf("Warning: failed to cache session: %v\n", err)
		}
	}

	if err := g.cache.UpdateLastUsed(sessionID); err != nil {
		fmt.Printf("Warning: failed to update last used: %v\n", err)
	}

	prompt := g.buildPrompt(diff, isSummarized)

	model := &opencode.Model{
		ProviderID: g.config.Generation.Model.Provider,
		ModelID:    g.config.Generation.Model.ModelID,
	}

	response, err := g.client.SendMessage(sessionID, prompt, model)
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	message := extractCommitMessage(response)
	return message, nil
}

/**
 * buildPrompt creates the AI prompt with diff and style instructions.
 *
 * @param diff - The git diff to include in the prompt
 * @param isSummarized - Whether the diff was summarized due to size
 * @returns The complete prompt string
 */
func (g *Generator) buildPrompt(diff string, isSummarized bool) string {
	style := g.config.Generation.Style
	styleGuide := getStyleGuide(style)

	var summarizedNote string
	if isSummarized {
		summarizedNote = `
NOTE: The diff below has been summarized because the original was too large.
Focus on the file list, diff stat, and available code changes to understand the overall changes.
`
	}

	prompt := fmt.Sprintf(`You are a git commit message generator. Your task is to generate a concise, meaningful commit message based on the following code changes.

%s
%s
Generate ONLY the commit message, nothing else. No explanation, no markdown formatting, just the message.

Here are the staged changes:

%s`, styleGuide, summarizedNote, diff)

	return prompt
}

/**
 * getStyleGuide returns the prompt instructions for the specified style.
 *
 * @param style - The commit style (conventional, imperative, detailed)
 * @returns The style guide instructions
 */
func getStyleGuide(style string) string {
	switch style {
	case "imperative":
		return `Follow the imperative mood style:
- Start with a verb (Add, Remove, Fix, Update, etc.)
- Write in the imperative mood, as if commanding someone
- Keep it under 72 characters
- Example: "Add user authentication to login page"`

	case "detailed":
		return `Use a detailed style with scope:
- Format: type(scope): description
- Scope should be short. dont write the whole file name. make it short and clear
- Types: feat, fix, docs, style, refactor, perf, test, chore
- Include a brief description in the body if needed
- Example: "feat(auth): add user authentication to login page
- Example if long filenames(eg. client_domain_person_check): "feat(domain): add user authentication to login page"`

	default:
		return `Follow the Conventional Commits style:
- Format: type(scope): description
- Scope should be short. dont write the whole file name. make it short and clear
- Types: feat, fix, docs, style, refactor, perf, test, chore
- Keep the description under 72 characters
- Example: "feat(auth): add user authentication
- Example if long filenames(eg. client_domain_person_check): "feat(domain): add user authentication to login page"`
	}
}

/**
 * extractCommitMessage extracts the clean commit message from AI response.
 *
 * @param response - The raw AI response
 * @returns The cleaned commit message (first line only)
 */
func extractCommitMessage(response string) string {
	response = strings.TrimSpace(response)

	if strings.HasPrefix(response, "```") {
		lines := strings.Split(response, "\n")
		if len(lines) > 1 {
			response = strings.Join(lines[1:], "\n")
		}
	}

	if before, ok := strings.CutSuffix(response, "```"); ok {
		response = before
	}

	lines := strings.Split(response, "\n")
	message := strings.TrimSpace(lines[0])

	return message
}
