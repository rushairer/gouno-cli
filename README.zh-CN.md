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
go mod tidy
make dev
# → http://localhost:8080
```

**参数：**

| 参数 | 简写 | 默认值 | 说明 |
|------|------|--------|------|
| `--module` | `-m` | 项目名称 | Go module 路径（如 `github.com/you/project`） |
| `--template` | `-t` | `./templates` | 本地路径或 Git URL 指向模板目录 |
| `--template-set` | | | 模板集名称（保存到 `.gouno.yaml` 供代码生成使用） |

**示例：**

```bash
# 使用默认模板
gouno-cli new my-api -m github.com/me/my-api

# 使用指定模板集
gouno-cli new order-service --template-set gorm -m github.com/me/order-service

# 使用自定义模板仓库
gouno-cli new my-app -t https://github.com/myorg/custom-template -m github.com/me/my-app

# 使用本地模板目录
gouno-cli new my-app -t /path/to/local/template -m github.com/me/my-app
```

### 管理模板集

```bash
# 列出已安装的模板集
gouno-cli template list

# 从 Git 安装模板集
gouno-cli template install gorm https://github.com/myorg/gouno-template-gorm

# 从本地路径安装
gouno-cli template install my-local /path/to/template

# 删除模板集
gouno-cli template remove gorm
```

模板集存储在 `~/.gouno/templates/` 目录下。

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
4. 失败时自动清理所有已创建的部分文件。

## 相关项目

| 仓库 | 说明 |
|------|------|
| [gouno](https://github.com/rushairer/gouno) | 核心库 |
| [gouno-template](https://github.com/rushairer/gouno-template) | 默认模板集 |
| [gouno-doc](https://github.com/rushairer/gouno-doc) | 文档 |

## 许可证

MIT 许可证。详见 [LICENSE](LICENSE)。
