# gouno-cli

[English](./README.md) | [文档](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/)

---

从 [gouno-template](https://github.com/rushairer/gouno-template) 脚手架生成 Go Web 项目的命令行工具。

## 安装

```bash
go install github.com/rushairer/gouno-cli@latest
```

或从源码构建：

```bash
git clone https://github.com/rushairer/gouno-cli
cd gouno-cli
go build -o gouno-cli .
```

## 使用方法

### 创建新项目

```bash
gouno-cli new my-service -m github.com/you/my-service
```

此命令会克隆默认的 [gouno-template](https://github.com/rushairer/gouno-template)，渲染所有模板变量，并创建一个可直接运行的项目。

```bash
cd my-service
make dev
# → http://localhost:8080
```

**参数：**

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--module` | `-m` | 项目名称 | Go module 路径（如 `github.com/you/project`） |
| `--template` | `-t` | `./templates` | 本地路径或 Git URL 指向模板目录 |
| `--skip-tidy` | | `false` | 创建项目后跳过 `go mod tidy` |

**示例：**

```bash
# 使用默认模板
gouno-cli new my-api -m github.com/me/my-api

# 使用自定义模板仓库
gouno-cli new my-app -t https://github.com/myorg/custom-template -m github.com/me/my-app

# 使用本地模板目录
gouno-cli new my-app -t /path/to/local/template -m github.com/me/my-app
```

### 关于项目模板

`gouno-cli new` 使用**项目模板** —— 完整的 Go 项目骨架仓库
（如 [gouno-template](https://github.com/rushairer/gouno-template)）。
可通过 `--template` 指定任意 Git URL 或本地目录，无需维护本地模板库。

> 注意：请勿将项目模板与**模板集（template set）**混淆。模板集是
> [gouno](https://github.com/rushairer/gouno) 库中 `gouno gen` 生成代码所用的
> `.tmpl` 脚手架文件集合，属于另一个概念，不属于 gouno-cli。

### 查看版本

```bash
gouno-cli version
# 或
gouno-cli --version
```

## 工作原理

1. `gouno-cli new` 克隆模板仓库（默认或指定的）。
2. 包含 `{{` 的文件使用提供的 module 路径和项目名进行 Go 模板渲染。
3. 其他文件直接复制（跳过 `.git/`、`templates/`、`bin/`）。
4. 默认自动执行 `go mod tidy`，除非传入 `--skip-tidy`。
5. 失败时自动清理所有已创建的部分文件。

## 相关项目

| 仓库 | 说明 |
|------|------|
| [gouno](https://github.com/rushairer/gouno) | 核心库（含内置 `.tmpl` 模板的代码生成） |
| [gouno-template](https://github.com/rushairer/gouno-template) | 默认项目模板 |
| [gouno-doc](https://github.com/rushairer/gouno-doc) | 文档 |

## 许可证

MIT 许可证。详见 [LICENSE](LICENSE)。
