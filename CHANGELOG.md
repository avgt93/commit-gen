# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.3] - 2026-02-20

### Added
- Interactive confirmation flow: users can now accept, edit, or regenerate commit messages before applying them.
- New `--no-confirm` (or `-n`) flag for the `generate` command to skip the interactive prompt.
- Added `generation.confirm` configuration option (default: `true`) to control the confirmation behavior.
- Added `git.editor` configuration option to specify a preferred text editor for message editing.

### Changed
- Updated default AI model to `opencode/gpt-5-nano` for improved performance and cost-efficiency.
- Improved git hook behavior: the hook now populates the commit message and leverages git's native editor flow for confirmation/editing.
- Refactored codebase for better clarity by removing redundant internal comments while maintaining comprehensive docstrings.

### Fixed
- Fixed a silent failure in the git hook when the configured model was unavailable.
- Fixed an issue where the git hook could fail in certain multi-branch environments by improving repository root detection and error handling.

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

[0.1.3]: https://github.com/avgt93/commit-gen/releases/tag/v0.1.3
[0.1.0]: https://github.com/avgt93/commit-gen/releases/tag/v0.1.0
