# Gouno Project Template Contract v1

> Status: **Stable / Normative**  
> Bootstrapper: `gouno-cli new`  
> Introduced as a documented contract: `gouno-cli v1.2.1`  
> Authority: `rushairer/gouno-cli`

This document defines the contract between `gouno-cli new` and a Gouno project template. It describes bootstrap mechanics only. It does not prescribe a web framework, application architecture, database, configuration library, or code-generator catalog.

The default `rushairer/gouno-template` is a reference implementation, not the definition of this contract.

## 1. Template source

`gouno-cli new` accepts:

- a local directory;
- an HTTPS Git repository URL;
- an SSH Git repository URL beginning with `git@`.

The default template source is `https://github.com/rushairer/gouno-template` when the historical local `./templates` default does not exist.

For a remote template, `--template-ref` is passed to a shallow Git clone as the branch/tag selector. If it is omitted, the repository's default branch is used. Reproducible project creation SHOULD pin an immutable tag or otherwise immutable ref.

## 2. Bootstrap inputs

The first-stage bootstrap renderer exposes exactly these template data fields:

- `{{.ModulePath}}`: the Go module path supplied with `--module` (`-m`), or the project name when omitted;
- `{{.ProjectName}}`: the project name supplied as the positional argument to `gouno-cli new`.

Bootstrap rendering applies to file contents, not file or directory names.

Project names and module paths are validated before template processing. A project is not allowed to overwrite an existing destination directory.

## 3. First-stage rendering

For ordinary template files:

1. files whose contents do not contain `{{` are copied verbatim;
2. files containing `{{` are parsed as Go `text/template` using the bootstrap data above;
3. if parsing succeeds, the template is executed and the rendered content is written;
4. if parsing fails, the file is copied verbatim rather than being treated as a bootstrap template;
5. if execution fails after a successful parse, project creation fails and rolls back.

Template authors SHOULD use bootstrap expressions only where project creation actually needs substitution. Do not rely on parse failure as a namespace mechanism; Codegen runtime resources have an explicit raw-copy rule described below.

## 4. Template-owned runtime resources

Codegen v1 uses a second-stage template language after the project has been created. Therefore these paths are copied **verbatim**, without first-stage bootstrap rendering:

```text
.gouno/codegen.yaml
.gouno/codegen/**
```

This rule preserves expressions such as:

```text
{{ arg "name" }}
{{ flag "path" }}
{{ camel ... }}
```

until `gouno gen` executes inside the generated project.

The normative Codegen protocol is defined by `rushairer/gouno` in `docs/codegen-template-spec.md`. `gouno-cli` does not define generator names or generator policy.

A template that does not support Codegen may omit `.gouno/codegen.yaml` and `.gouno/codegen/` entirely.

## 5. Reserved and filtered paths

During bootstrap, any path containing one of these path segments is skipped:

```text
.git
.idea
.DS_Store
bin
templates
.env
```

A path segment beginning with `.env.` is also skipped, as is any path segment ending in `.local.yaml`.

### Legacy `templates/` reservation

The `templates/` filtering rule predates Codegen v1 and is retained for compatibility with older template repositories. It is a legacy reserved path, not the current location for code-generation templates.

New Codegen-enabled templates MUST place runtime generator resources under `.gouno/codegen/`, not under `templates/`.

Removal or reinterpretation of this legacy reservation is a compatibility-sensitive contract change and must not happen silently.

## 6. File modes

Directories created in the generated project use mode `0755`.

Copied/rendered files preserve the source file's permission bits subject to a `0755` mask. This preserves executable files while preventing group/other write bits from being propagated by the bootstrapper.

## 7. Module tidying

Unless `--skip-tidy` is supplied, `gouno-cli new` runs:

```bash
go mod tidy
```

inside the generated project after copying/rendering completes.

A template intended for normal use SHOULD therefore produce a project whose module files can be tidied successfully in a clean environment.

## 8. Failure and rollback

If template copying/rendering fails, or if the automatic `go mod tidy` step fails, the partially created destination directory is removed.

Templates SHOULD be validated by creating a temporary project from the same source/ref that users will consume and then running the generated project's own quality checks.

## 9. Two-stage template model

Gouno templates can contain two distinct template languages with different lifetimes:

```text
Stage 1: project bootstrap
  owner: gouno-cli
  time:  gouno-cli new
  data:  .ModulePath / .ProjectName

Stage 2: project code generation (optional)
  owner: gouno runtime protocol + project template policy
  time:  gouno gen ...
  data:  args / flags / Codegen v1 functions
```

`.gouno/codegen.yaml` and `.gouno/codegen/**` form the explicit boundary between these stages and are raw-copied during Stage 1.

## 10. Versioning and compatibility

Template authors SHOULD publish immutable tags for versions intended for reuse and document which Gouno runtime versions they require.

A template version can change its own architecture and generator policy because those are template concerns. However, it must continue to satisfy the bootstrap contract of the `gouno-cli` versions it claims to support and the Codegen schema version it declares.

Changes to this document that alter existing bootstrap behavior require compatibility review and release notes. New behavior must not be inferred from the current default template alone.
