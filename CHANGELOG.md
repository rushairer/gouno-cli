# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [Unreleased]

### Changed

- Preserve `.gouno/codegen.yaml` and `.gouno/codegen/**` verbatim when creating projects so template-owned runtime codegen expressions are not consumed by the project-bootstrap renderer.
- Keep project-template rendering and runtime codegen rendering as separate stages with separate ownership.

## [1.2.0] - 2026-08-24

### Added

- Add `--template-ref` for pinning an immutable remote template branch or tag.

### Changed

- Require Go 1.25.0 or newer and validate Go 1.25.x and 1.26.x in CI.
- `gouno-cli new` now follows the template repository's default branch by default (including the default `gouno-template`) with a shallow clone; use `--template-ref` to pin an immutable branch or tag for reproducible builds.

### Security

- Add SHA-pinned shared quality and Conventional PR CI gates plus Dependabot.

## [1.1.1] - 2026-08-20

### Fixed

- Skip private environment and local config files (`*.local.yaml`, `.env`, `.env.*`) when copying project template skeleton (`gouno/new.go`).
- Fall back to reading module version from runtime build info (`debug.ReadBuildInfo`) when `Version` is not injected via `-ldflags` (`gouno/root.go`).

## [1.1.0] - 2026-08-20

### Breaking

- Remove the `template install/list/remove` commands. `gouno-cli new` uses project templates (full skeleton repositories like gouno-template) directly via `--template <git-url-or-path>`; there is no local template registry to maintain.
- Remove the `--template-set` flag from `new`; no longer writes `.gouno.yaml`.

### Changed

- Clarify the distinction between a project template (used by `new`) and a template set (`.tmpl` scaffold files for code generation in the gouno library).
- Refuse to overwrite an existing project directory.

### Fixed

- Stop wrongly skipping `.gitignore` and `.github` files caused by the `.git` prefix check (now matches the exact name).
- Validate each path element of the module path, rejecting empty, `.` and `..` segments.
- Keep executable bits on generated files while stripping group/other write permissions, preventing overly permissive modes (e.g. 0777) on scaffolded projects.

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
