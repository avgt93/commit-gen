# commit-gen

A CLI tool that generates descriptive commit messages using OpenCode's AI capabilities. Simply run `git commit -m ""` and let AI analyze your staged changes to create meaningful commit messages.

## Features

- AI-powered commit message generation using OpenCode
- Dual mode support: run via CLI subprocess or connect to OpenCode server
- Automatic large diff summarization (handles diffs > 32KB)
- Multiple commit styles: conventional, imperative, detailed
- Session caching for faster subsequent commits (server mode)
- Git hook integration for automatic message generation
- Highly configurable via YAML, environment variables, or CLI flags

## Installation

### Prerequisites

- Go 1.21 or later (for building from source)
- OpenCode installed and available in PATH
- Git repository

### From Source

```bash
git clone https://github.com/avgt93/commit-gen
cd commit-gen
make install
```

### Available Make Commands

```bash
make build       # Build the binary
make install     # Build and install to /usr/local/bin
make test        # Run all tests
make clean       # Remove build artifacts
make run         # Build and run the CLI
make lint        # Run linter
make fmt         # Format code
make release     # Build cross-platform releases
```

### Manual Build

```bash
go build -o commit-gen ./cmd/commit-gen
sudo mv commit-gen /usr/local/bin/
```

## Quick Start

### 1. Verify OpenCode is available

```bash
commit-gen health
```

### 2. Install the git hook in your repository

```bash
cd /path/to/your/repo
commit-gen install
```

### 3. Make changes and stage them

```bash
git add .
```

### 4. Commit with empty message to trigger AI generation

```bash
git commit
```

The tool will analyze your staged changes and generate a descriptive commit message automatically.

## Usage

### Commands

```
Available Commands:
  cache       Manage session cache
  completion  Generate the autocompletion script for the specified shell
  config      Manage configuration
  generate    Generate a commit message from staged changes
  health      Check if the OpenCode backend is available
  help        Help about any command
  init        Initialize the configuration file
  install     Install git hook for automatic commit message generation
  preview     Preview changes and generated commit message
  uninstall   Remove the git hook
  version     Show version information
```

### Generate a Commit Message

```bash
# Generate and apply commit message
commit-gen generate

# Preview without applying
commit-gen generate --dry-run

# Specify commit style
commit-gen generate --style imperative

# Use server mode instead of default run mode
commit-gen generate --mode server
```

### Preview Changes

```bash
# Show staged diff and generated message
commit-gen preview
```

### Configuration Management

```bash
# View current configuration
commit-gen config

# Initialize config file
commit-gen init
```

### Cache Management (Server Mode)

```bash
# Show cache status
commit-gen cache status

# Clear all cached sessions
commit-gen cache clear
```

### Git Hook Management

```bash
# Install hook
commit-gen install

# Remove hook
commit-gen uninstall
```

### Health Check

```bash
# Check OpenCode backend availability
commit-gen health
```

## Operation Modes

### Run Mode (Default)

Uses `opencode run` CLI command directly. No server required.

```bash
commit-gen generate --mode run
```

Benefits:
- No need to start OpenCode server
- Simpler setup
- Faster for single commits

### Server Mode

Connects to OpenCode HTTP API server.

```bash
# Start server in another terminal
opencode serve

# Use server mode
commit-gen generate --mode server
```

Benefits:
- Session caching for context reuse
- Better for frequent commits
- Supports concurrent requests

## Configuration

Configuration hierarchy (highest to lowest priority):
1. CLI flags
2. Environment variables (`COMMIT_GEN_*` prefix)
3. Config file (`~/.config/commit-gen/config.yaml`)
4. Default values

### Config File

Create `~/.config/commit-gen/config.yaml`:

```yaml
opencode:
  mode: run              # "run" or "server"
  host: localhost        # server mode only
  port: 4096             # server mode only
  timeout: 120

generation:
  style: conventional    # conventional, imperative, detailed
  model:
    provider: google
    model_id: antigravity-gemini-3-pro

cache:
  enabled: true          # server mode only
  ttl: 24h

git:
  staged_only: true
  editor: cat
  max_diff_size: 32768   # bytes before summarizing (32KB default)
```

### Environment Variables

```bash
export COMMIT_GEN_OPENCODE_MODE=run
export COMMIT_GEN_OPENCODE_HOST=localhost
export COMMIT_GEN_OPENCODE_PORT=4096
export COMMIT_GEN_GENERATION_STYLE=conventional
export COMMIT_GEN_GENERATION_MODEL_PROVIDER=google
export COMMIT_COMMIT_GEN_GENERATION_MODEL_MODEL_ID=antigravity-gemini-3-pro
export COMMIT_GEN_GIT_MAX_DIFF_SIZE=32768
```

## Commit Styles

### Conventional (Default)

Format: `type(scope): description`

Types: feat, fix, docs, style, refactor, perf, test, chore

Examples:
- `feat(auth): add user authentication`
- `fix(api): handle null pointer exception`
- `docs(readme): update installation steps`

### Imperative

Uses the imperative mood, as if commanding someone.

Examples:
- `Add user authentication to login page`
- `Fix null pointer exception in API handler`
- `Update README with installation steps`

### Detailed

Format: `type(scope): description` with optional body

Examples:
- `feat(auth): add user authentication`
- `fix(api): handle null pointer exception in getUser endpoint`

## Large Diff Handling

When staged changes exceed 32KB (configurable via `git.max_diff_size`), the diff is automatically summarized for AI processing. The summary includes:

- List of changed files
- Diff statistics (insertions/deletions)
- Truncated diff content
- Note to AI about summarization

This prevents failures with large commits while still providing meaningful context.

## Troubleshooting

### "opencode binary not found in PATH"

Ensure OpenCode is installed and available:
```bash
which opencode
```

### "opencode server is not running" (Server Mode)

Start OpenCode server:
```bash
opencode serve
```

### "no staged changes found"

Stage your changes first:
```bash
git add .
```

### Hook not working

Verify installation:
```bash
cat .git/hooks/prepare-commit-msg
```

Reinstall if needed:
```bash
commit-gen uninstall
commit-gen install
```

### Permission denied when installing

Run with sudo:
```bash
sudo make install
```

Or manually copy binary:
```bash
make build
sudo cp commit-gen /usr/local/bin/
```

## Development

### Project Structure

```
commit-gen/
├── cmd/commit-gen/          # CLI entry point
├── internal/
│   ├── git/                 # Git operations
│   ├── opencode/            # OpenCode client and runner
│   ├── config/              # Configuration management
│   ├── cache/               # Session caching
│   ├── generator/           # Commit message generation
│   └── hook/                # Git hook management
├── Makefile
└── README.md
```

### Build and Test

```bash
# Build
make build

# Run tests
make test

# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean
```

See [AGENTS.md](./AGENTS.md) for detailed architecture documentation.

## Contributing

Contributions welcome:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT

## Support

- GitHub Issues: [submit an issue](https://github.com/avgt93/commit-gen/issues)
- OpenCode Discord: [join the community](https://opencode.ai/discord)

## Acknowledgments

- Built with [OpenCode](https://opencode.ai)
- CLI powered by [Cobra](https://github.com/spf13/cobra)
- Config with [Viper](https://github.com/spf13/viper)
