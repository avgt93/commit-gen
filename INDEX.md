# commit-gen - Complete Index

A Go CLI tool that generates commit messages using OpenCode's AI. Just run `git commit -m ""` and it fills in the message!

## 📁 Project Files (16 total)

### Documentation (5 files)
| File | Purpose |
|------|---------|
| **README.md** | Main user guide, installation, usage |
| **GETTING_STARTED.md** | Quick start tutorial with examples |
| **AGENTS.md** | Architecture & design documentation |
| **PROJECT_SUMMARY.md** | Project overview & quick reference |
| **STRUCTURE.txt** | Visual project structure & flow |

### Source Code (7 files)
| File | Lines | Purpose |
|------|-------|---------|
| **cmd/commit-gen/main.go** | 262 | CLI entry point, all commands |
| **internal/opencode/client.go** | 179 | HTTP client for OpenCode server |
| **internal/cache/session_cache.go** | 174 | Session caching with TTL |
| **internal/generator/commit.go** | 161 | Core generation logic |
| **internal/git/diff.go** | 120 | Git operations |
| **internal/hook/install.go** | 125 | Git hook management |
| **internal/config/config.go** | 104 | Configuration handling |

### Configuration (4 files)
| File | Purpose |
|------|---------|
| **go.mod** | Go module dependencies |
| **go.sum** | Dependency checksums |
| **Makefile** | Build automation |
| **.gitignore** | Git ignore patterns |

### Executable
| File | Size | Purpose |
|------|------|---------|
| **commit-gen** | 11 MB | Compiled binary (ready to use!) |

---

## 🚀 Quick Start

### 1. Start OpenCode Server
```bash
opencode serve
```

### 2. Install Hook in Your Repo
```bash
cd /path/to/your/repo
/home/avgt/all/kanban/commit-gen/commit-gen install
```

### 3. Use It!
```bash
git add .
git commit -m ""
# AI generates commit message automatically!
```

See **GETTING_STARTED.md** for detailed walkthrough.

---

## 📖 Documentation Guide

### For First-Time Users
1. Start with **README.md** - Overview & installation
2. Follow **GETTING_STARTED.md** - Step-by-step tutorial
3. Try the commands!

### For Understanding Architecture
1. Read **AGENTS.md** - Detailed architecture & design
2. Check **STRUCTURE.txt** - Visual structure & execution flow
3. Read the source code in `internal/`

### For Project Overview
- **PROJECT_SUMMARY.md** - Quick reference with statistics
- **STRUCTURE.txt** - File-by-file breakdown

---

## 🛠️ Key Features

✅ **AI-Powered**: Uses OpenCode's Claude 3.5 Sonnet  
✅ **Git Integration**: Hooks into git automatically  
✅ **Smart Caching**: Reuses sessions for speed  
✅ **Multiple Styles**: Conventional, imperative, detailed  
✅ **Configuration**: YAML + environment variables  
✅ **Single Binary**: No external dependencies  

---

## 📝 CLI Commands

```bash
# Show help
commit-gen --help

# Generate a message
commit-gen generate
commit-gen generate --style imperative
commit-gen generate --dry-run

# Preview before committing
commit-gen preview

# Manage hook
commit-gen install
commit-gen uninstall

# View config
commit-gen config

# Cache management
commit-gen cache status
commit-gen cache clear

# Version
commit-gen version
```

---

## 🔧 Source Code Organization

```
internal/
├── git/           - Git operations (diff, messages, status)
├── opencode/      - HTTP client to OpenCode server
├── config/        - Configuration management (Viper)
├── cache/         - Session caching with MD5 keys & TTL
├── generator/     - Commit message generation logic
└── hook/          - Git hook installation/management
```

**Total Lines of Code**: 1,125 lines (clean, modular Go)

---

## ⚙️ Configuration

### Config File
`~/.config/commit-gen/config.yaml`:
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

### Environment Variables
```bash
export COMMIT_GEN_GENERATION_STYLE=imperative
export COMMIT_GEN_OPENCODE_HOST=localhost
export COMMIT_GEN_OPENCODE_PORT=4096
```

---

## 🏗️ How It Works

1. User runs: `git commit -m ""`
2. Git triggers: `prepare-commit-msg` hook
3. Hook runs: `commit-gen generate --hook`
4. Tool checks: OpenCode server running
5. Gets: Staged git diff
6. Sends to OpenCode: Diff + AI prompt
7. Receives: Generated commit message
8. Writes: Message to `.git/COMMIT_EDITMSG`
9. Git completes: Commit with AI-generated message

See **STRUCTURE.txt** for visual diagram.

---

## 🚨 Troubleshooting

### OpenCode server not running?
```bash
opencode serve
```

### No staged changes?
```bash
git add .
```

### Hook not working?
```bash
commit-gen uninstall
commit-gen install
```

See **GETTING_STARTED.md** for more troubleshooting.

---

## 📚 File Quick Reference

| Want to... | Read... |
|-----------|---------|
| Get started quickly | **GETTING_STARTED.md** |
| Install & setup | **README.md** |
| Understand architecture | **AGENTS.md** |
| See project overview | **PROJECT_SUMMARY.md** |
| View file structure | **STRUCTURE.txt** |
| Build from source | **Makefile** |
| Check config options | **internal/config/config.go** |
| See CLI commands | **cmd/commit-gen/main.go** |

---

## 🔨 Build & Development

### Build
```bash
cd /home/avgt/all/kanban/commit-gen
go build -o commit-gen ./cmd/commit-gen
```

### Run
```bash
./commit-gen --help
./commit-gen generate
```

### Install to /usr/local/bin
```bash
make install
```

### Clean
```bash
make clean
```

---

## 📦 Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/spf13/viper` - Configuration
- `github.com/fatih/color` - Terminal colors
- Go standard library for HTTP, git, file operations

No external system dependencies - single standalone binary!

---

## 💡 Key Design Decisions

1. **Shell Commands for Git**: Simple, reliable, no extra dependencies
2. **Session Caching**: Speeds up repeated commits in same repo
3. **Persistent Cache**: `~/.cache/commit-gen/sessions.json`
4. **HTTP Client**: Direct OpenCode API, no SDK needed
5. **Modular Structure**: Each package has single responsibility
6. **Cobra CLI**: Industry-standard Go CLI framework
7. **Viper Config**: Flexible configuration (files + env vars)

---

## 🎯 Next Steps

1. **Test It**: Run on a real repository with staged changes
2. **Customize**: Edit config or commit styles as needed
3. **Share**: Binary is portable and self-contained
4. **Contribute**: Extend with new features if desired

---

## 📄 File Sizes

- **Binary**: 11 MB (standalone executable)
- **Total Code**: 1,125 lines of Go
- **Documentation**: ~20 KB of markdown
- **Config**: YAML-based, minimal

---

## ✨ What Makes This Special

✅ Single binary - no installation needed  
✅ Uses OpenCode AI - high-quality messages  
✅ Session caching - fast after first use  
✅ Git hook integration - seamless workflow  
✅ Well-documented - multiple guides included  
✅ Clean code - modular, easy to understand  
✅ Configurable - styles, timeouts, models  
✅ Error handling - clear messages when things go wrong  

---

## 📞 Support

For issues:
- Check **GETTING_STARTED.md** troubleshooting
- Read **AGENTS.md** for architecture details
- Check OpenCode docs: https://opencode.ai/docs

---

**Last Updated**: February 2, 2026  
**Version**: 0.1.0  
**Location**: `/home/avgt/all/kanban/commit-gen`
