# commit-gen Project Summary

## ✅ Project Successfully Created!

All components have been built and tested. Here's what's included:

### Project Structure

```
/home/avgt/all/kanban/commit-gen/
├── cmd/commit-gen/
│   └── main.go                 # CLI entry point with all commands
├── internal/
│   ├── git/
│   │   └── diff.go             # Git operations (diff, status, messages)
│   ├── opencode/
│   │   └── client.go           # HTTP client for OpenCode server
│   ├── config/
│   │   └── config.go           # Viper-based config management
│   ├── cache/
│   │   └── session_cache.go    # Session caching with TTL
│   ├── generator/
│   │   └── commit.go           # Core generation logic
│   └── hook/
│       └── install.go          # Git hook installation
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
├── Makefile                    # Build targets
├── README.md                   # User documentation
├── AGENTS.md                   # Architecture documentation
├── .gitignore                  # Git ignore patterns
└── commit-gen                  # Compiled binary (ready to use!)
```

## Features Implemented

✅ **Git Integration**
- Get staged changes via `git diff --staged`
- Read/write commit message files
- Get repository information

✅ **OpenCode Client**
- HTTP connection to OpenCode server
- Session creation and management
- Message sending with model selection

✅ **Session Caching**
- In-memory + persistent caching
- MD5 hashing of repo paths as keys
- 24-hour TTL (configurable)
- Cache status and clear commands

✅ **Commit Message Generation**
- Three styles: conventional, imperative, detailed
- Smart prompt generation
- Response parsing and cleanup

✅ **Git Hook Integration**
- Auto-installation of `prepare-commit-msg` hook
- Empty message detection
- Message file writing

✅ **CLI Commands**
- `generate` - Generate commit message
- `install` - Install git hook
- `uninstall` - Remove git hook
- `config` - Show configuration
- `preview` - Preview changes + message
- `cache status` - Cache info
- `cache clear` - Clear cache
- `version` - Show version

✅ **Configuration**
- YAML-based configuration
- Environment variable overrides (`COMMIT_GEN_*`)
- Sensible defaults

## Getting Started

### 1. Start OpenCode Server
```bash
opencode serve
```

### 2. Install Git Hook in Your Repository
```bash
cd /path/to/your/repo
/home/avgt/all/kanban/commit-gen/commit-gen install
```

### 3. Stage Your Changes
```bash
git add .
```

### 4. Commit with Empty Message to Trigger AI Generation
```bash
git commit -m ""
```

The tool will automatically generate a descriptive commit message!

## Quick Commands

```bash
# Show help
./commit-gen --help

# Generate a message manually
./commit-gen generate

# Preview changes and generated message
./commit-gen preview

# Specific commit style
./commit-gen generate --style imperative

# Show configuration
./commit-gen config

# Cache management
./commit-gen cache status
./commit-gen cache clear

# Install/Uninstall hook
./commit-gen install
./commit-gen uninstall
```

## Building from Source

```bash
cd /home/avgt/all/kanban/commit-gen

# Build
make build

# Run
./commit-gen --help

# Install to /usr/local/bin (requires sudo)
make install

# Test
make test
```

## Configuration File

Create `~/.config/commit-gen/config.yaml`:

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

git:
  staged_only: true
```

## How It Works

1. **User runs**: `git commit -m ""`
2. **Git triggers**: `prepare-commit-msg` hook
3. **Hook runs**: `commit-gen generate --hook`
4. **Tool checks**: OpenCode server is running
5. **Tool retrieves**: Staged git diff
6. **Tool sends**: Diff + prompt to OpenCode
7. **OpenCode returns**: Generated commit message
8. **Tool writes**: Message to `.git/COMMIT_EDITMSG`
9. **Git uses**: Message for the commit

## Key Files Overview

### main.go (271 lines)
- Cobra CLI framework setup
- All command definitions
- Configuration initialization
- Color-coded output

### client.go (160 lines)
- HTTP client for OpenCode API
- Session management
- Message sending and response handling

### session_cache.go (180 lines)
- In-memory caching
- Persistent cache files
- TTL-based expiration
- MD5 repo path hashing

### commit.go (140 lines)
- Generation orchestration
- OpenCode integration
- Style-based prompts
- Response parsing

### config.go (100 lines)
- Configuration loading
- Viper setup
- Environment variable support
- Config accessors

### diff.go (110 lines)
- Git diff retrieval
- File status checking
- Repository information
- Message file operations

### install.go (100 lines)
- Git hook script generation
- Hook installation/removal
- Hook validation

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration
- `github.com/fatih/color` - Colored output
- Go standard library (net/http, os/exec, etc.)

All as Go modules - no external system dependencies!

## Next Steps

1. **Test with a repository**: Try running `git commit -m ""` on a repo
2. **Customize prompts**: Edit styles in `generator/commit.go` if needed
3. **Add to PATH**: `export PATH=/home/avgt/all/kanban/commit-gen:$PATH`
4. **Share with team**: Binary is self-contained and portable

## Troubleshooting

**"opencode server is not running"**
- Run `opencode serve` in another terminal
- Make sure it's on localhost:4096

**"no staged changes found"**
- Stage your changes: `git add .`

**Hook not working**
- Verify: `cat .git/hooks/prepare-commit-msg`
- Reinstall: `commit-gen uninstall && commit-gen install`

## Project Documentation

- **README.md** - User guide and installation
- **AGENTS.md** - Architecture and design documentation
- **Makefile** - Build automation

Enjoy! 🚀
