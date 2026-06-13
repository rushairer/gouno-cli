# gouno-cli

[中文](./README.zh-CN.md) | [Documentation](https://github.com/rushairer/gouno-doc)

---

A CLI tool to scaffold Go web projects from [gouno-template](https://github.com/rushairer/gouno-template).

## Install

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

This clones the default [gouno-template](https://github.com/rushairer/gouno-template), renders all template variables, and creates a ready-to-run project.

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
| `--template-set` | | | Template set name (saved to `.gouno.yaml` for code generation) |
| `--skip-tidy` | | `false` | Skip running `go mod tidy` after project creation |

**Examples:**

```bash
# Use default template
gouno-cli new my-api -m github.com/me/my-api

# Use a specific template set
gouno-cli new order-service --template-set gorm -m github.com/me/order-service

# Use a custom template repository
gouno-cli new my-app -t https://github.com/myorg/custom-template -m github.com/me/my-app

# Use a local template directory
gouno-cli new my-app -t /path/to/local/template -m github.com/me/my-app
```

### Manage Template Sets

```bash
# List installed template sets
gouno-cli template list

# Install a template set from git
gouno-cli template install gorm https://github.com/myorg/gouno-template-gorm

# Install from local path
gouno-cli template install my-local /path/to/template

# Remove a template set
gouno-cli template remove gorm
```

Template sets are stored in `~/.gouno/templates/`.

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
| [gouno](https://github.com/rushairer/gouno) | Core library |
| [gouno-template](https://github.com/rushairer/gouno-template) | Default template set |
| [gouno-doc](https://github.com/rushairer/gouno-doc) | Documentation |

## License

MIT License. See [LICENSE](LICENSE) for details.
