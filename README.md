# gouno-cli

[中文](./README.zh-CN.md) | [Documentation](https://github.com/rushairer/gouno-doc)

---

A CLI tool to scaffold Go web projects from [gouno-template](https://github.com/rushairer/gouno-template).

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

### Create a New Project

```bash
gouno-cli new my-service -m github.com/you/my-service
```

This clones the default [gouno-template](https://github.com/rushairer/gouno-template) at immutable tag `v1.2.0`, renders all template variables, and creates a ready-to-run project.

```bash
cd my-service
make dev
# → http://localhost:8080
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--module` | `-m` | project name | Go module path (e.g., `github.com/you/project`) |
| `--template` | `-t` | `./templates` | Local path or git URL to template directory |
| `--template-ref` | | empty | Immutable branch, tag, or commit ref for a remote template |
| `--skip-tidy` | | `false` | Skip running `go mod tidy` after project creation |

**Examples:**

```bash
# Use default template
gouno-cli new my-api -m github.com/me/my-api

# Use a custom template repository
gouno-cli new my-app -t https://github.com/myorg/custom-template --template-ref v1.0.0 -m github.com/me/my-app

# Use a local template directory
gouno-cli new my-app -t /path/to/local/template -m github.com/me/my-app
```

### About Project Templates

`gouno-cli new` uses a project template — a full Go project skeleton repository
(e.g. [gouno-template](https://github.com/rushairer/gouno-template)). You can
point it to any git URL or local directory with `--template`; there is no local
template registry to maintain.

For reproducible remote builds, provide an immutable `--template-ref`. Custom remote templates without it retain their repository default-branch behavior.

> Note: don't confuse a project template with a *template set*. A template set
> is the collection of `.tmpl` scaffold files used by `gouno gen` (from the
> [gouno](https://github.com/rushairer/gouno) library) to generate code — it is
> a different concept and not part of gouno-cli.

### Version

```bash
gouno-cli version
# or
gouno-cli --version
```

## How It Works

1. `gouno-cli new` clones a template repository (default or specified).
2. Files containing `{{` are rendered as Go templates using the provided module path and project name.
3. Other files are copied as-is (skipping `.git/`, `templates/`, `bin/`).
4. `go mod tidy` runs automatically unless `--skip-tidy` is set.
5. On failure, all partially created files are cleaned up automatically.

## Related Projects

| Repository | Description |
|------------|-------------|
| [gouno](https://github.com/rushairer/gouno) | Core library (includes code generation with built-in `.tmpl` templates) |
| [gouno-template](https://github.com/rushairer/gouno-template) | Default project template |
| [gouno-doc](https://github.com/rushairer/gouno-doc) | Documentation |

## License

MIT License. See [LICENSE](LICENSE) for details.
