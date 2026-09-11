# Gouno CLI Agent Contract

This file defines repository-maintenance invariants for humans and coding agents working on `rushairer/gouno-cli`.

## Product boundary

`gouno-cli` is an architecture-agnostic project bootstrapper. It creates projects from full project-template repositories or local directories.

It does not own application architecture, framework choice, or code-generator policy. Do not introduce assumptions about DDD, Clean Architecture, Gin, Cobra, Viper, database technology, or generator names merely because the default `gouno-template` currently uses them.

## Normative bootstrap contract

The normative contract between `gouno-cli new` and a project template is `docs/project-template-contract.md`.

Changes to any of the following are contract changes and must be treated as compatibility-sensitive:

- accepted template sources or `--template-ref` behavior;
- bootstrap variables or rendering semantics;
- ignored/reserved paths;
- file permission handling;
- `.gouno` runtime-resource preservation;
- `go mod tidy` behavior;
- rollback/error behavior.

When changing these semantics, update the contract, tests, user-facing README when applicable, and `CHANGELOG.md` under `Unreleased` in the same change.

## Codegen boundary

Codegen policy belongs to the project template, not this CLI.

For Codegen v1, `.gouno/codegen.yaml` and `.gouno/codegen/**` are template-owned runtime resources and must be copied verbatim during project bootstrap. They contain a second-stage Go-template language that must not be consumed by the first-stage bootstrap renderer.

Do not restore the historical model in which `gouno-cli` or Gouno Core owns a generator template set.

## Historical compatibility

The top-level path segment `templates/` is still filtered by `gouno-cli new` for historical compatibility. New templates must not use it for Codegen v1 resources; use `.gouno/codegen/` instead.

Do not remove or reinterpret a legacy filtering rule casually. First assess existing-template compatibility and make any deprecation/removal explicit in the Project Template Contract and release notes.

## Cross-repository boundary

Related responsibilities are intentionally split:

- `gouno-cli`: project-template bootstrap contract and project creation;
- `gouno`: reusable mechanisms and Codegen protocol/runtime;
- `gouno-template`: default template/reference implementation;
- `gouno-doc`: user and template-author guides.

Do not duplicate the normative Codegen schema here. Link to the authoritative specification in `rushairer/gouno`.

## Validation

Behavioral changes to project creation must have tests that exercise template copying/rendering and rollback semantics. Changes affecting `.gouno` resources must verify that Codegen expressions survive bootstrap unchanged.

Follow all repository CI, lint, vet, race, vulnerability, and release requirements before merging.
