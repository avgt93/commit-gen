# commit-gen

A CLI tool that generates descriptive commit messages using OpenCode's AI capabilities.

## Features

- 🤖 AI-powered commit message generation using OpenCode
- ⚡ Fast session caching to reuse context
- 📝 Multiple commit styles (conventional, imperative, detailed)
- 🔗 Simple git hook integration
- ⚙️ Highly configurable
- 🎯 Just run `git commit -m ""` and let AI fill it in!

## Installation

### From source

```bash
git clone https://github.com/avgt93/commit-gen
cd commit-gen
go build -o commit-gen ./cmd/commit-gen
sudo mv commit-gen /usr/local/bin/
```

## Quick Start

### 1. Start OpenCode server (in a separate terminal)

```bash
opencode serve
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
git commit -m ""
```

The tool will automatically generate a descriptive commit message!

## Usage

### Generate a commit message manually

```bash
commit-gen generate
```

### Preview changes and generated message

```bash
commit-gen preview
```

### Specify commit style

```bash
commit-gen generate --style imperative
```

Styles: `conventional`, `imperative`, `detailed`

### View configuration

```bash
commit-gen config
```

### Manage cache

```bash
# Show cache status
commit-gen cache status

# Clear cache
commit-gen cache clear
```

### Uninstall git hook

```bash
commit-gen uninstall
```

## Configuration

Create `~/.config/commit-gen/config.yaml`:

```yaml
opencode:
  host: localhost
  port: 4096
  timeout: 30

generation:
  style: conventional
  model:
    provider: anthropic
    model_id: claude-3-5-sonnet-20241022

cache:
  enabled: true
  ttl: 24h
```

Or use environment variables:

```bash
export COMMIT_GEN_OPENCODE_HOST=localhost
export COMMIT_GEN_OPENCODE_PORT=4096
export COMMIT_GEN_GENERATION_STYLE=conventional
```

## How It Works

1. **Detects empty commit message** - Git hook intercepts `git commit -m ""`
2. **Retrieves staged changes** - Gets the diff of staged files
3. **Connects to OpenCode** - Sends diff to OpenCode server
4. **Generates message** - AI generates a commit message based on changes
5. **Applies message** - Message is written to the commit

## Commit Styles

### Conventional Commits
Format: `type(scope): description`

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

## Requirements

- Go 1.21 or later (only if building from source)
- OpenCode server running (`opencode serve`)
- Git repository

## Troubleshooting

### "opencode server is not running"

Make sure OpenCode is running in another terminal:
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

## Development

### Build

```bash
go build -o commit-gen ./cmd/commit-gen
```

### Test

```bash
go test ./...
```

### Run

```bash
./commit-gen --help
```

## Architecture

See [AGENTS.md](./AGENTS.md) for detailed architecture documentation.

## Contributing

Contributions welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

MIT

## Support

For issues and questions:
- GitHub Issues: [submit an issue](https://github.com/avgt93/commit-gen/issues)
- OpenCode Discord: [join the community](https://opencode.ai/discord)

## Acknowledgments

- Built with [OpenCode](https://opencode.ai)
- CLI powered by [Cobra](https://github.com/spf13/cobra)
- Config with [Viper](https://github.com/spf13/viper)
