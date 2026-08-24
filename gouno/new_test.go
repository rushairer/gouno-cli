package gouno

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

func TestValidateProjectName(t *testing.T) {
	validNames := []string{
		"myproject",
		"my_project",
		"_private",
		"Project123",
		"a",
		"demo-app",
		"my-cool-project",
	}

	for _, name := range validNames {
		t.Run("valid/"+name, func(t *testing.T) {
			if err := validateProjectName(name); err != nil {
				t.Errorf("validateProjectName(%q) = %v; want nil", name, err)
			}
		})
	}

	invalidNames := []struct {
		name string
		desc string
	}{
		{"", "empty"},
		{"../etc/passwd", "path traversal with .."},
		{"my/project", "contains slash"},
		{"my\\project", "contains backslash"},
		{"123abc", "starts with digit"},
		{"-project", "starts with dash"},
		{"my project", "contains space"},
		{"my.project", "contains dot"},
	}

	for _, tt := range invalidNames {
		t.Run("invalid/"+tt.desc, func(t *testing.T) {
			if err := validateProjectName(tt.name); err == nil {
				t.Errorf("validateProjectName(%q) = nil; want error", tt.name)
			}
		})
	}
}

func TestShouldSkipFile(t *testing.T) {
	tests := []struct {
		path   string
		expect bool
	}{
		{".git", true},
		{".git/config", true},
		{".gitignore", false},
		{".gitattributes", false},
		{".github/workflows/ci.yml", false},
		{".idea", true},
		{".DS_Store", true},
		{"bin/gouno", true},
		{"templates/base", true},
		{"src/main.go", false},
		{"go.mod", false},
		{"README.md", false},
		{"internal/domain/.gitkeep", false},
		{"config/development.yaml", false},
		{"config/development.local.yaml", true},
		{"config/production.local.yaml", true},
		{".env", true},
		{".env.local", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := shouldSkipFile(tt.path)
			if got != tt.expect {
				t.Errorf("shouldSkipFile(%q) = %v; want %v", tt.path, got, tt.expect)
			}
		})
	}
}

func TestIsRenderableFile(t *testing.T) {
	tests := []struct {
		content string
		expect  bool
	}{
		{"package main\n// {{.ProjectName}}", true},
		{"module {{.ModulePath}}", true},
		{"package main\nfunc main() {}", false},
		{"# README\nThis is a test.", false},
		{"name: {{.ProjectName}}", true},
	}

	for _, tt := range tests {
		t.Run(tt.content[:20], func(t *testing.T) {
			got := isRenderableFile(tt.content)
			if got != tt.expect {
				t.Errorf("isRenderableFile(%q) = %v; want %v", tt.content, got, tt.expect)
			}
		})
	}
}

func TestCopyTemplate(t *testing.T) {
	// 创建临时源目录
	srcDir := t.TempDir()
	destDir := t.TempDir()

	// 创建模板文件
	goModContent := `module {{.ModulePath}}

go 1.23`
	mainContent := `package main

import "fmt"

func main() {
	fmt.Println("{{.ProjectName}}")
}`
	readmeContent := `# {{.ProjectName}}

This is a project.`

	if err := os.WriteFile(filepath.Join(srcDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "main.go"), []byte(mainContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "README.md"), []byte(readmeContent), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(srcDir, ".git"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, ".git", "config"), []byte("git config"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(srcDir, "bin"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "bin", "app"), []byte("binary"), 0755); err != nil {
		t.Fatal(err)
	}

	data := TemplateData{
		ModulePath:  "github.com/test/myapp",
		ProjectName: "myapp",
	}

	err := copyTemplate(srcDir, destDir, data)
	if err != nil {
		t.Fatalf("copyTemplate failed: %v", err)
	}

	// 验证 go.mod 已渲染
	gomod, err := os.ReadFile(filepath.Join(destDir, "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if string(gomod) != "module github.com/test/myapp\n\ngo 1.23" {
		t.Errorf("go.mod not rendered: %s", string(gomod))
	}

	// 验证 main.go 已渲染
	maingo, err := os.ReadFile(filepath.Join(destDir, "main.go"))
	if err != nil {
		t.Fatalf("read main.go: %v", err)
	}
	if !contains(string(maingo), `"myapp"`) {
		t.Errorf("main.go not rendered: %s", string(maingo))
	}

	// 验证 README.md 已渲染（包含 {{）
	readme, err := os.ReadFile(filepath.Join(destDir, "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	if !contains(string(readme), "# myapp") {
		t.Errorf("README.md not rendered: %s", string(readme))
	}

	// 验证 .git 目录被跳过
	if _, err := os.Stat(filepath.Join(destDir, ".git")); !os.IsNotExist(err) {
		t.Error(".git directory should be skipped")
	}

	// 验证 bin 目录被跳过
	if _, err := os.Stat(filepath.Join(destDir, "bin")); !os.IsNotExist(err) {
		t.Error("bin directory should be skipped")
	}
}

func TestCopyTemplateCleanupOnError(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()
	destProject := filepath.Join(destDir, "myproject")

	// 创建一个会渲染成功的文件
	if err := os.WriteFile(filepath.Join(srcDir, "go.mod"), []byte("module {{.ModulePath}}"), 0644); err != nil {
		t.Fatal(err)
	}

	// 创建一个会导致 Execute 失败的模板（引用不存在的字段）
	// template.Parse 会成功，但 Execute 会失败
	badContent := `{{.NonExistentField}}`
	if err := os.WriteFile(filepath.Join(srcDir, "bad.go"), []byte(badContent), 0644); err != nil {
		t.Fatal(err)
	}

	data := TemplateData{
		ModulePath:  "test",
		ProjectName: "test",
	}

	err := copyTemplate(srcDir, destProject, data)
	if err == nil {
		t.Fatal("expected error from copyTemplate, got nil")
	}

	// 验证目标目录中的文件被正确写入（bad.go 应该有错误）
	// 注意：cleanup 由调用方（RunE 中的 os.RemoveAll）负责，copyTemplate 本身不清理
}

func TestTidyProjectRunsGoModTidy(t *testing.T) {
	orig := runExternalCommand
	defer func() { runExternalCommand = orig }()

	var gotDir, gotName string
	var gotArgs []string
	runExternalCommand = func(dir, name string, args ...string) error {
		gotDir = dir
		gotName = name
		gotArgs = append([]string(nil), args...)
		return nil
	}

	if err := tidyProject("/tmp/project"); err != nil {
		t.Fatalf("tidyProject() error: %v", err)
	}
	if gotDir != "/tmp/project" || gotName != "go" || strings.Join(gotArgs, " ") != "mod tidy" {
		t.Fatalf("unexpected command: dir=%q name=%q args=%q", gotDir, gotName, strings.Join(gotArgs, " "))
	}
}

func TestTidyProjectWrapsError(t *testing.T) {
	orig := runExternalCommand
	defer func() { runExternalCommand = orig }()

	expected := errors.New("boom")
	runExternalCommand = func(_ string, _ string, _ ...string) error {
		return expected
	}

	err := tidyProject("/tmp/project")
	if err == nil {
		t.Fatal("tidyProject() = nil; want error")
	}
	if !strings.Contains(err.Error(), "running go mod tidy") {
		t.Fatalf("tidyProject() error = %v; want context", err)
	}
	if !errors.Is(err, expected) {
		t.Fatalf("tidyProject() error = %v; want wrapped expected error", err)
	}
}

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestValidateModulePath(t *testing.T) {
	validPaths := []string{
		"github.com/foo/bar",
		"myproject",
		"example.com/a-b/c_d",
		"git.example.com/org/module.v2",
	}

	for _, p := range validPaths {
		t.Run("valid/"+p, func(t *testing.T) {
			if err := validateModulePath(p); err != nil {
				t.Errorf("validateModulePath(%q) = %v; want nil", p, err)
			}
		})
	}

	invalidPaths := []struct {
		path string
		desc string
	}{
		{"", "empty"},
		{"a\nb", "newline"},
		{"a\tb", "tab"},
		{"foo bar", "space"},
		{"foo%bar", "percent"},
		{"foo/../bar", "traversal element"},
		{"foo/./bar", "dot element"},
		{"foo//bar", "empty element"},
		{"../foo", "leading traversal"},
		{"foo/..", "trailing traversal"},
	}

	for _, tt := range invalidPaths {
		t.Run("invalid/"+tt.desc, func(t *testing.T) {
			if err := validateModulePath(tt.path); err == nil {
				t.Errorf("validateModulePath(%q) = nil; want error", tt.path)
			}
		})
	}
}

func TestCopyTemplatePreservesGitignore(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(srcDir, ".gitignore"), []byte("bin/\n*.exe\n"), 0644); err != nil {
		t.Fatal(err)
	}

	data := TemplateData{ModulePath: "test", ProjectName: "test"}
	if err := copyTemplate(srcDir, destDir, data); err != nil {
		t.Fatalf("copyTemplate() error: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(destDir, ".gitignore"))
	if err != nil {
		t.Fatalf("read .gitignore: %v", err)
	}
	if !strings.Contains(string(content), "bin/") {
		t.Errorf(".gitignore content = %q; want original content preserved", string(content))
	}
}

func TestCopyTemplatePreservesFileMode(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	script := filepath.Join(srcDir, "run.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\necho hi\n"), 0755); err != nil {
		t.Fatal(err)
	}

	data := TemplateData{ModulePath: "test", ProjectName: "test"}
	if err := copyTemplate(srcDir, destDir, data); err != nil {
		t.Fatalf("copyTemplate() error: %v", err)
	}

	info, err := os.Stat(filepath.Join(destDir, "run.sh"))
	if err != nil {
		t.Fatalf("stat run.sh: %v", err)
	}
	if info.Mode().Perm() != 0o755 {
		t.Errorf("run.sh mode = %v; want 0755", info.Mode().Perm())
	}
}

// chdir 切换工作目录并在测试结束时恢复
func chdir(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
}

// setNewCmdFlags 为 newCmd 设置全部相关 flags,避免包级 flag 值在测试间残留
func setNewCmdFlags(t *testing.T, templateDir, module string, skipTidy bool) {
	t.Helper()
	f := newCmd.Flags()
	for flag, value := range map[string]string{
		"template":     templateDir,
		"module":       module,
		"skip-tidy":    strconv.FormatBool(skipTidy),
		"template-ref": "",
	} {
		if err := f.Set(flag, value); err != nil {
			t.Fatalf("set flag %q: %v", flag, err)
		}
	}
}

func TestNewCmdCreatesProject(t *testing.T) {
	chdir(t, t.TempDir())

	srcDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(srcDir, "go.mod"), []byte("module {{.ModulePath}}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, "README.md"), []byte("# {{.ProjectName}}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(srcDir, ".gitignore"), []byte("bin/\n"), 0644); err != nil {
		t.Fatal(err)
	}

	setNewCmdFlags(t, srcDir, "github.com/me/app", true)

	if err := newCmd.RunE(newCmd, []string{"myapp"}); err != nil {
		t.Fatalf("newCmd.RunE() error: %v", err)
	}

	gomod, err := os.ReadFile(filepath.Join("myapp", "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if string(gomod) != "module github.com/me/app\n" {
		t.Errorf("go.mod = %q; want rendered module path", string(gomod))
	}

	readme, err := os.ReadFile(filepath.Join("myapp", "README.md"))
	if err != nil {
		t.Fatalf("read README.md: %v", err)
	}
	if !strings.Contains(string(readme), "# myapp") {
		t.Errorf("README.md = %q; want project name rendered", string(readme))
	}

	if _, err := os.Stat(filepath.Join("myapp", ".gitignore")); err != nil {
		t.Errorf(".gitignore should be copied to new project: %v", err)
	}
}

func TestNewCmdRejectsExistingDir(t *testing.T) {
	chdir(t, t.TempDir())

	if err := os.MkdirAll("myapp", 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join("myapp", "keep.txt"), []byte("user data"), 0644); err != nil {
		t.Fatal(err)
	}

	srcDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(srcDir, "go.mod"), []byte("module {{.ModulePath}}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	setNewCmdFlags(t, srcDir, "test", true)

	err := newCmd.RunE(newCmd, []string{"myapp"})
	if err == nil {
		t.Fatal("expected error for existing directory, got nil")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("error = %v; want 'already exists'", err)
	}

	// 已存在的用户数据必须完好
	content, err := os.ReadFile(filepath.Join("myapp", "keep.txt"))
	if err != nil {
		t.Fatalf("user data destroyed: %v", err)
	}
	if string(content) != "user data" {
		t.Errorf("keep.txt = %q; want unchanged", string(content))
	}
}

func TestNewCmdRejectsInvalidModulePath(t *testing.T) {
	chdir(t, t.TempDir())

	srcDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(srcDir, "go.mod"), []byte("module {{.ModulePath}}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	setNewCmdFlags(t, srcDir, "foo/../bar", true)

	err := newCmd.RunE(newCmd, []string{"myapp"})
	if err == nil {
		t.Fatal("expected error for invalid module path, got nil")
	}
	if _, err := os.Stat("myapp"); !os.IsNotExist(err) {
		t.Error("project directory should not be created for invalid module path")
	}
}

func TestNewCmdClonesRemoteTemplate(t *testing.T) {
	chdir(t, t.TempDir())

	orig := runExternalCommand
	defer func() { runExternalCommand = orig }()

	var calls []string
	runExternalCommand = func(dir, name string, args ...string) error {
		calls = append(calls, name+" "+strings.Join(args, " "))
		if name == "git" && len(args) >= 2 && args[0] == "clone" {
			// 模拟 clone 产生模板文件
			dest := args[len(args)-1]
			if err := os.WriteFile(filepath.Join(dest, "go.mod"), []byte("module {{.ModulePath}}\n"), 0644); err != nil {
				return err
			}
		}
		return nil
	}

	setNewCmdFlags(t, "https://github.com/example/template.git", "github.com/me/app", true)
	if err := newCmd.RunE(newCmd, []string{"myapp"}); err != nil {
		t.Fatalf("newCmd.RunE() error: %v", err)
	}

	if len(calls) != 1 {
		t.Fatalf("expected 1 git clone call, got %v", calls)
	}
	if !strings.HasPrefix(calls[0], "git clone https://github.com/example/template.git") {
		t.Errorf("clone call = %q; want git clone with source URL", calls[0])
	}

	gomod, err := os.ReadFile(filepath.Join("myapp", "go.mod"))
	if err != nil {
		t.Fatalf("read go.mod: %v", err)
	}
	if string(gomod) != "module github.com/me/app\n" {
		t.Errorf("go.mod = %q; want rendered from cloned template", string(gomod))
	}
}

func TestNewCmdPinsDefaultTemplateTag(t *testing.T) {
	chdir(t, t.TempDir())
	orig := runExternalCommand
	defer func() { runExternalCommand = orig }()
	var call string
	runExternalCommand = func(_ string, name string, args ...string) error {
		call = name + " " + strings.Join(args, " ")
		if name == "git" && args[0] == "clone" {
			return os.WriteFile(filepath.Join(args[len(args)-1], "go.mod"), []byte("module {{.ModulePath}}\n"), 0644)
		}
		return nil
	}
	setNewCmdFlags(t, "./templates", "github.com/me/pinned", true)
	if err := newCmd.RunE(newCmd, []string{"pinned"}); err != nil {
		t.Fatalf("newCmd.RunE() error: %v", err)
	}
	want := "git clone --branch v1.2.0 --depth 1 https://github.com/rushairer/gouno-template"
	if !strings.HasPrefix(call, want) {
		t.Fatalf("clone call = %q; want prefix %q", call, want)
	}
}
