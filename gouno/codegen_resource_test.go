package gouno

import (
	"os"
	"path/filepath"
	"testing"
)

func TestTemplateRuntimeResourceDetection(t *testing.T) {
	for _, path := range []string{
		".gouno/codegen.yaml",
		".gouno/codegen/service.tmpl",
		filepath.Join(".gouno", "codegen", "nested", "service.tmpl"),
	} {
		if !isTemplateRuntimeResource(path) {
			t.Fatalf("expected %q to be a template runtime resource", path)
		}
	}
	if isTemplateRuntimeResource("templates/service.tmpl") {
		t.Fatal("legacy project template directory must not be treated as runtime codegen resources")
	}
}

func TestCopyTemplatePreservesCodegenTemplatesVerbatim(t *testing.T) {
	src := t.TempDir()
	dst := filepath.Join(t.TempDir(), "project")
	manifest := `schema: gouno.dev/codegen/v1
command:
  use: gen
generators:
  - name: service
    outputs:
      - template: .gouno/codegen/service.tmpl
        path: '{{ .Args.name }}.go'
`
	tmpl := `package service

type {{ .Args.name }}Service struct{}
`
	for path, content := range map[string]string{
		filepath.Join(src, ".gouno", "codegen.yaml"):            manifest,
		filepath.Join(src, ".gouno", "codegen", "service.tmpl"): tmpl,
	} {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	if err := copyTemplate(src, dst, TemplateData{ModulePath: "example.com/app", ProjectName: "app"}); err != nil {
		t.Fatal(err)
	}
	for rel, want := range map[string]string{
		filepath.Join(".gouno", "codegen.yaml"):            manifest,
		filepath.Join(".gouno", "codegen", "service.tmpl"): tmpl,
	} {
		got, err := os.ReadFile(filepath.Join(dst, rel))
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != want {
			t.Fatalf("runtime codegen resource %s was rendered unexpectedly\nwant:\n%s\ngot:\n%s", rel, want, got)
		}
	}
}
