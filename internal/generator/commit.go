package generator

import (
	"fmt"
	"strings"

	"github.com/avgt93/commit-gen/internal/cache"
	"github.com/avgt93/commit-gen/internal/config"
	"github.com/avgt93/commit-gen/internal/git"
	"github.com/avgt93/commit-gen/internal/opencode"
)

type Generator struct {
	client *opencode.Client
	cache  *cache.SessionCache
	config *config.Config
}

func NewGenerator(cfg *config.Config, cacheInstance *cache.SessionCache) *Generator {
	client := opencode.NewClient(cfg.OpenCode.Host, cfg.OpenCode.Port, cfg.OpenCode.Timeout)
	return &Generator{
		client: client,
		cache:  cacheInstance,
		config: cfg,
	}
}

func (g *Generator) Generate() (string, error) {
	healthy, err := g.client.CheckHealth()
	if err != nil || !healthy {
		return "", fmt.Errorf("opencode server is not running at %s:%d\n\nPlease start it with: opencode serve", g.config.OpenCode.Host, g.config.OpenCode.Port)
	}

	diff, err := git.GetStagedDiff()
	if err != nil {
		return "", fmt.Errorf("failed to get git diff: %w", err)
	}

	if strings.TrimSpace(diff) == "" {
		return "", fmt.Errorf("no staged changes found")
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

	g.cache.UpdateLastUsed(sessionID)

	// Build the prompt
	prompt := g.buildPrompt(diff)

	// Get model configuration
	model := &opencode.Model{
		ProviderID: g.config.Generation.Model.Provider,
		ModelID:    g.config.Generation.Model.ModelID,
	}

	// Send to OpenCode and get response
	response, err := g.client.SendMessage(sessionID, prompt, model)
	if err != nil {
		return "", fmt.Errorf("failed to generate commit message: %w", err)
	}

	// Extract just the commit message from the response
	message := extractCommitMessage(response)

	return message, nil
}

// buildPrompt constructs the prompt for OpenCode
func (g *Generator) buildPrompt(diff string) string {
	style := g.config.Generation.Style

	styleGuide := getStyleGuide(style)

	prompt := fmt.Sprintf(`You are a git commit message generator. Your task is to generate a concise, meaningful commit message based on the following code changes.

%s

Generate ONLY the commit message, nothing else. No explanation, no markdown formatting, just the message.

Here are the staged changes:

%s`, styleGuide, diff)

	return prompt
}

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
- Types: feat, fix, docs, style, refactor, perf, test, chore
- Include a brief description in the body if needed
- Example: "feat(auth): add user authentication to login page"`

	default: // conventional
		return `Follow the Conventional Commits style:
- Format: type(scope): description
- Types: feat, fix, docs, style, refactor, perf, test, chore
- Keep the description under 72 characters
- Example: "feat(auth): add user authentication"`
	}
}

func extractCommitMessage(response string) string {
	// Clean up the response - remove markdown code blocks if present
	response = strings.TrimSpace(response)

	// Remove markdown code blocks
	if strings.HasPrefix(response, "```") {
		lines := strings.Split(response, "\n")
		if len(lines) > 1 {
			response = strings.Join(lines[1:], "\n")
		}
	}

	if strings.HasSuffix(response, "```") {
		response = strings.TrimSuffix(response, "```")
	}

	// Take only the first line if multiple lines are returned
	lines := strings.Split(response, "\n")
	message := strings.TrimSpace(lines[0])

	return message
}
