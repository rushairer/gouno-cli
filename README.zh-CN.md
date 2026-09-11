# gouno-cli

[English](./README.md) | [文档](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/)

---

一个**与项目架构无关**的 Go 项目脚手架 CLI，可从完整的项目模板仓库或本地目录创建项目。

`gouno-cli` 负责项目 Bootstrap 机制，但**不规定**应用架构、Web 框架或代码生成策略。官方 [gouno-template](https://github.com/rushairer/gouno-template) 是默认参考模板，并不代表所有 Gouno 项目都必须采用同样的结构。

## 安装

需要 Go 1.25.0 或更高版本。CI 验证 Go 1.25.x 和 1.26.x；不再支持 Go 1.23/1.24。

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

在默认配置下，如果当前目录存在 `./templates`，CLI 会使用它；否则会克隆官方 [gouno-template](https://github.com/rushairer/gouno-template) 的默认分支。

```bash
cd my-service
make dev
```

**参数：**

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--module` | `-m` | 项目名称 | Go module 路径，例如 `github.com/you/project` |
| `--template` | `-t` | `./templates` | 本地模板目录或受支持的 Git 仓库 URL |
| `--template-ref` | | 空 | 远程模板的 branch/tag 选择器；为空时跟随其默认分支 |
| `--skip-tidy` | | `false` | 创建项目后跳过 `go mod tidy` |

**示例：**

```bash
# 使用官方默认模板
gouno-cli new my-api -m github.com/me/my-api

# 使用自定义远程模板
gouno-cli new my-app \
  -t https://github.com/myorg/custom-gouno-template \
  -m github.com/me/my-app

# 固定已发布模板版本，保证脚手架过程可复现
gouno-cli new my-app \
  -t https://github.com/myorg/custom-gouno-template \
  --template-ref v1.0.0 \
  -m github.com/me/my-app

# 使用本地模板目录
gouno-cli new my-app \
  -t /path/to/local/template \
  -m github.com/me/my-app
```

需要可复现创建项目时，`--template-ref` 应优先固定到不可变的 Release Tag，而不是持续移动的分支。

## Project Template 与 Codegen

**Project Template（项目模板）**是 `gouno-cli new` 使用的完整项目骨架。只要满足 Bootstrap Contract，它可以自行决定目录结构和技术栈。

Codegen 是可选能力，它属于 Template / 当前项目，不属于 `gouno-cli`，也不再是 Gouno Core 内置的一套 DDD Generator 清单。

支持 Codegen 的 Template 可以携带：

```text
.gouno/
├── codegen.yaml
└── codegen/
    └── ...
```

Codegen v1 资源在项目 Bootstrap 阶段会被**原样复制**，从而保留其中第二阶段使用的模板表达式。如果 Template 不提供 `.gouno/codegen.yaml`，生成项目的 CLI 完全可以没有 `gen` 命令。

项目创建行为的唯一规范见 [Project Template Contract v1](./docs/project-template-contract.md)。Codegen Schema 本身由 [gouno](https://github.com/rushairer/gouno/blob/main/docs/codegen-template-spec.md) 维护。

如何开发自己的 Project Template，参见 [Gouno 文档：Template Authoring](https://github.com/rushairer/gouno-doc/blob/main/zh-CN/template-authoring.md)。

## 两阶段模板模型

Gouno 项目可以同时存在两套生命周期完全不同的渲染：

```text
Stage 1: gouno-cli new
  项目 Bootstrap
  {{.ModulePath}} / {{.ProjectName}}

Stage 2: gouno gen ...（可选）
  项目自己的 Codegen
  args / flags / Codegen v1 模板函数
```

`.gouno/codegen.yaml` 与 `.gouno/codegen/**` 是这两个阶段之间明确的 raw-copy 边界。

## 项目创建流程

1. 克隆或读取所选的完整 Project Template。
2. 对普通 Bootstrap 模板文件按需使用 `ModulePath` 与 `ProjectName` 渲染内容。
3. 将 `.gouno/codegen*` 下的 Codegen v1 Runtime Resources 原样复制。
4. 跳过 `.git/`、`bin/`、`.env*`、`*.local.yaml` 以及历史保留的 `templates/` 等路径。
5. 在安全权限掩码下保留可执行位。
6. 默认运行 `go mod tidy`，除非使用 `--skip-tidy`。
7. 渲染或依赖整理失败时删除未完成的目标项目。

`templates/` 的过滤属于历史兼容规则；新的 Codegen 资源应放在 `.gouno/codegen/`。

## 查看版本

```bash
gouno-cli version
# 或
gouno-cli --version
```

## 相关项目

| 仓库 | 说明 |
|------|------|
| [gouno](https://github.com/rushairer/gouno) | 可复用机制以及 Codegen 协议/运行时 |
| [gouno-template](https://github.com/rushairer/gouno-template) | 官方默认 Project Template 与参考 Codegen Policy |
| [gouno-doc](https://github.com/rushairer/gouno-doc) | 用户、Template Authoring 与工程规范文档 |

## 许可证

MIT 许可证。详见 [LICENSE](LICENSE)。
