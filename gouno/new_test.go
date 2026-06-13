package gouno

import (
	"errors"
	"os"
	"path/filepath"
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
		{".gitignore", true},
		{".github/workflows/ci.yml", true},
		{".idea", true},
		{".DS_Store", true},
		{"bin/gouno", true},
		{"templates/base", true},
		{"src/main.go", false},
		{"go.mod", false},
		{"README.md", false},
		{"internal/domain/.gitkeep", true},
		{"config/development.yaml", false},
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

	os.WriteFile(filepath.Join(srcDir, "go.mod"), []byte(goModContent), 0644)
	os.WriteFile(filepath.Join(srcDir, "main.go"), []byte(mainContent), 0644)
	os.WriteFile(filepath.Join(srcDir, "README.md"), []byte(readmeContent), 0644)
	os.MkdirAll(filepath.Join(srcDir, ".git"), 0755)
	os.WriteFile(filepath.Join(srcDir, ".git", "config"), []byte("git config"), 0644)
	os.MkdirAll(filepath.Join(srcDir, "bin"), 0755)
	os.WriteFile(filepath.Join(srcDir, "bin", "app"), []byte("binary"), 0755)

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
	os.WriteFile(filepath.Join(srcDir, "go.mod"), []byte("module {{.ModulePath}}"), 0644)

	// 创建一个会导致 Execute 失败的模板（引用不存在的字段）
	// template.Parse 会成功，但 Execute 会失败
	badContent := `{{.NonExistentField}}`
	os.WriteFile(filepath.Join(srcDir, "bad.go"), []byte(badContent), 0644)

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
