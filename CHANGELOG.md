# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.0] - 2026-02-02

### Added
- Initial release of commit-gen
- Core message generation with OpenCode integration
- Git hooks support (prepare-commit-msg automatic hook)
- Session caching with TTL (24-hour default) and persistence
- YAML configuration support with environment variable overrides
- Three commit message styles: conventional, imperative, detailed
- CLI commands: generate, install, uninstall, config, preview, cache status, cache clear, version
- Comprehensive test suite (43+ unit tests)
- Complete documentation and architecture guides
- Support for multiple platforms (Linux, macOS, Windows)

### Features
- **Message Generation**: Integrates with OpenCode server for AI-powered commit messages
- **Git Integration**: Automatically retrieves staged diffs and manages commit messages
- **Session Management**: Caches sessions to improve performance
- **Configuration Management**: Flexible YAML config + environment variables
- **Hook Installation**: Easy setup with `commit-gen install`
- **Multiple Styles**: Choose between conventional, imperative, or detailed commit formats

### Known Limitations
- Requires OpenCode server to be running locally (default: localhost:4096)
- Git repository tests require actual git environment
- Integration tests require both git repo and OpenCode server

[0.1.0]: https://github.com/avgt93/commit-gen/releases/tag/v0.1.0
