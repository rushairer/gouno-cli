# gouno-cli

[中文](./README.zh-CN.md) | [Documentation](https://github.com/rushairer/gouno-doc)

---

An architecture-agnostic CLI for scaffolding Go projects from full project-template repositories or local directories.

`gouno-cli` owns project bootstrap mechanics. It does **not** prescribe application architecture, framework choice, or code-generator policy. The official [gouno-template](https://github.com/rushairer/gouno-template) is the default reference template, not the definition of what a Gouno project must look like.

## Install

Requires Go 1.25.0 or newer. CI validates Go 1.25.x and 1.26.x; Go 1.23/1.24 are unsupported.

```bash
go install github.com/rushairer/gouno-cli@latest
```

Or build from source:

```bash
git clone https://github.com/rushairer/gouno-cli
cd gouno-cli
go build -o gouno-cli .
```

## Usage

### Create a new project

```bash
gouno-cli new my-service -m github.com/you/my-service
```

With the default settings, `gouno-cli` uses a local `./templates` directory when one exists; otherwise it clones the official [gouno-template](https://github.com/rushairer/gouno-template) default branch.

```bash
cd my-service
make dev
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--module` | `-m` | project name | Go module path (for example `github.com/you/project`) |
| `--template` | `-t` | `./templates` | Local template directory or supported Git repository URL |
| `--template-ref` | | empty | Branch/tag selector for a remote template; empty follows its default branch |
| `--skip-tidy` | | `false` | Skip `go mod tidy` after project creation |

**Examples:**

```bash
# Official default template
gouno-cli new my-api -m github.com/me/my-api

# Custom remote template
gouno-cli new my-app \
  -t https://github.com/myorg/custom-gouno-template \
  -m github.com/me/my-app

# Pin a released template for reproducible scaffolding
gouno-cli new my-app \
  -t https://github.com/myorg/custom-gouno-template \
  --template-ref v1.0.0 \
  -m github.com/me/my-app

# Local template directory
gouno-cli new my-app \
  -t /path/to/local/template \
  -m github.com/me/my-app
```

For repeatable project creation, prefer an immutable release tag for `--template-ref` rather than a moving branch.

## Project templates and Codegen

A **project template** is a complete project skeleton consumed by `gouno-cli new`. It can choose any project structure and technology stack that satisfies the bootstrap contract.

Code generation is optional and belongs to the template/project, not to `gouno-cli` and not to a built-in DDD catalog in Gouno Core.

A Codegen-enabled template can ship:

```text
.gouno/
├── codegen.yaml
└── codegen/
    └── ...
```

These Codegen v1 resources are copied verbatim during project bootstrap so their second-stage template expressions remain intact. If a template does not provide `.gouno/codegen.yaml`, the generated project's CLI does not need to expose a `gen` command at all.

See the normative [Project Template Contract v1](./docs/project-template-contract.md) for bootstrap behavior. The Codegen schema itself is owned by [gouno](https://github.com/rushairer/gouno/blob/main/docs/codegen-template-spec.md).

For a user-facing guide to creating custom templates, see [Gouno Documentation](https://github.com/rushairer/gouno-doc/blob/main/template-authoring.md).

## Two-stage template model

Gouno projects may use two separate rendering stages:

```text
Stage 1: gouno-cli new
  Project bootstrap
  {{.ModulePath}} / {{.ProjectName}}

Stage 2: gouno gen ... (optional)
  Project-owned Codegen
  args / flags / Codegen v1 template functions
```

`.gouno/codegen.yaml` and `.gouno/codegen/**` form an explicit raw-copy boundary between those stages.

## How project bootstrap works

1. Clone/read the selected full project template.
2. Render ordinary bootstrap-template contents with `ModulePath` and `ProjectName` where applicable.
3. Copy Codegen v1 runtime resources under `.gouno/codegen*` verbatim.
4. Skip reserved/private paths such as `.git/`, `bin/`, `.env*`, `*.local.yaml`, and the legacy reserved `templates/` path.
5. Preserve executable permissions subject to the bootstrap safety mask.
6. Run `go mod tidy` unless `--skip-tidy` is set.
7. Remove the partial destination if rendering or module tidying fails.

The historical `templates/` filtering rule remains for compatibility; new Codegen resources belong under `.gouno/codegen/`.

## Version

```bash
gouno-cli version
# or
gouno-cli --version
```

## Related Projects

| Repository | Description |
|------------|-------------|
| [gouno](https://github.com/rushairer/gouno) | Reusable mechanisms plus the Codegen protocol/runtime |
| [gouno-template](https://github.com/rushairer/gouno-template) | Official default project template and reference Codegen policy |
| [gouno-doc](https://github.com/rushairer/gouno-doc) | User, template-authoring, and engineering documentation |

## License

MIT License. See [LICENSE](LICENSE) for details.
