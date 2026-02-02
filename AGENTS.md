# commit-gen: AI-Powered Commit Message Generator

## Project Overview

`commit-gen` is a CLI tool that automatically generates descriptive commit messages using OpenCode's AI capabilities. Instead of manually typing commit messages, users can simply run `git commit -m ""` and the tool will fill in the message based on their staged changes.

## Architecture

The project is structured into several modular packages:

```
commit-gen/
├── cmd/commit-gen/          # CLI entry point using Cobra
├── internal/
│   ├── git/                 # Git operations (diff, status, file management)
│   ├── opencode/            # HTTP client to OpenCode server
│   ├── config/              # Configuration management with Viper
│   ├── cache/               # Session caching mechanism
│   ├── generator/           # Core commit message generation logic
│   └── hook/                # Git hook installation/management
```

## How It Works

### 1. Git Integration
- When user runs `git commit -m ""`, the git hook is triggered
- The hook detects the empty message and invokes `commit-gen generate --hook`
- Git operations in `internal/git/diff.go` retrieve staged changes

### 2. OpenCode Server Connection
- Checks if OpenCode server is running at configured host:port (default: localhost:4096)
- If not running, prompts user to start it with: `opencode serve`

### 3. Session Management
- Creates a session in OpenCode for the current repository
- Sessions are cached based on repository path (MD5 hash)
- Cache TTL is configurable (default: 24 hours)
- Reuses sessions to avoid creating new ones for each commit

### 4. Message Generation
- Sends staged diff + system prompt to OpenCode
- OpenCode AI generates a descriptive commit message
- Currently supports three styles: conventional, imperative, detailed
- Extracts clean message from AI response (removes markdown, takes first line)

### 5. Git Hook Flow
- `prepare-commit-msg` hook intercepts empty messages
- Writes generated message to `.git/COMMIT_EDITMSG`
- User still has option to edit before committing

## Configuration

Configuration hierarchy (highest to lowest priority):
1. CLI flags
2. Environment variables (`COMMIT_GEN_*` prefix)
3. Config file (`~/.config/commit-gen/config.yaml`)
4. Default values

### Default Configuration
```yaml
opencode:
  host: localhost
  port: 4096
  timeout: 30

generation:
  style: conventional  # conventional, imperative, detailed
  model:
    provider: anthropic
    model_id: claude-3-5-sonnet-20241022

cache:
  enabled: true
  ttl: 24h
  location: ~/.cache/commit-gen

git:
  staged_only: true
```

## CLI Commands

### `commit-gen generate`
- Generates a commit message from staged changes
- Flags:
  - `--style`: Message style (conventional, imperative, detailed)
  - `--dry-run`: Preview without writing
  - `--hook`: Internal flag for hook usage

### `commit-gen install`
- Installs the `prepare-commit-msg` git hook in current repo

### `commit-gen uninstall`
- Removes the git hook

### `commit-gen config`
- Shows current configuration

### `commit-gen preview`
- Shows staged diff + generated message

### `commit-gen cache status`
- Shows cache statistics

### `commit-gen cache clear`
- Clears all cached sessions

## System Prompts

The tool uses intelligent prompts based on commit style:

### Conventional
```
Follow the Conventional Commits style:
- Format: type(scope): description
- Types: feat, fix, docs, style, refactor, perf, test, chore
- Keep the description under 72 characters
```

### Imperative
```
Follow the imperative mood style:
- Start with a verb (Add, Remove, Fix, Update, etc.)
- Write in the imperative mood, as if commanding someone
- Keep it under 72 characters
```

### Detailed
```
Use a detailed style with scope:
- Format: type(scope): description
- Include a brief description in the body if needed
```

## Session Caching Strategy

- **In-memory cache**: Active sessions during CLI execution
- **Persistent cache**: `~/.cache/commit-gen/sessions.json`
- **Cache key**: MD5 hash of repository path
- **TTL**: 24 hours (configurable)
- Benefits:
  - Avoids creating new OpenCode sessions for each commit
  - Faster message generation after first commit
  - Reuses context understanding

## Error Handling

### OpenCode Server Not Running
Prompts user with clear instruction:
```
Error: opencode server is not running at localhost:4096

To fix this, run:
  opencode serve

In another terminal, then try again.
```

### No Staged Changes
Returns error: "no staged changes found"

### Empty Response
Returns error: "no text response received"

## Extension Points

Future enhancements can be added at:

1. **New Commit Styles**: Add new cases in `getStyleGuide()` in `generator/commit.go`
2. **Custom Prompts**: Add prompt customization in configuration
3. **Model Selection**: Allow CLI-flag based model switching
4. **Post-processing**: Add filters in `extractCommitMessage()` for custom cleanup
5. **Integration**: Add support for GitHub commits, co-authored commits, etc.

## Development Notes

### Key Dependencies
- `github.com/spf13/cobra`: CLI framework
- `github.com/spf13/viper`: Configuration management
- `github.com/fatih/color`: Colored terminal output
- No external git binary dependency - uses shell commands via `os/exec`

### Package Responsibilities

**git/**: 
- Get staged diff
- Read/write commit messages
- Get repository root and name
- Manage commit message files

**opencode/**:
- HTTP client for OpenCode API
- Session creation and management
- Message sending and response handling

**config/**:
- Load and parse YAML configuration
- Handle environment variable overrides
- Provide configuration accessors

**cache/**:
- In-memory and persistent session caching
- TTL-based expiration
- Repository-based cache keys

**generator/**:
- Orchestrate the generation pipeline
- Build prompts based on style
- Parse AI responses

**hook/**:
- Install/uninstall git hooks
- Generate hook script content
- Validate hook installation

## Testing

Run existing tests:
```bash
go test ./...
```

Test individual packages:
```bash
go test ./internal/git
go test ./internal/cache
```

## Building and Distribution

Build binary:
```bash
go build -o commit-gen ./cmd/commit-gen
```

This produces a single, standalone executable with no dependencies.

