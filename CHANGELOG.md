# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Fixed

- Use secure default file permissions (0755 for directories, 0644 for files) instead of copying source template permissions directly, preventing overly permissive access control on scaffolded projects.

## [1.0.1] - 2026-06-13

### Changed

- Run `go mod tidy` automatically after project creation.
- Add `--skip-tidy` for offline or custom dependency workflows.
- Update English and Chinese usage docs for the streamlined new-project flow.

## [1.0.0] - 2026-05-31

### Added

- `new` command to scaffold projects from template sets.
- `template install` / `list` / `remove` commands for template set management.
- `--module`, `--template`, `--template-set` flags for `new` command.
- `--version` / `-v` flag and `version` subcommand.
- Automatic project name validation.
- Template rendering with Go's `text/template`.
- Cleanup on partial creation failure.
- README with installation guide and usage documentation.
