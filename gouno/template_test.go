package gouno

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/cobra"
)

func TestTemplateSetDir(t *testing.T) {
	dir, err := templateSetDir()
	if err != nil {
		t.Fatalf("templateSetDir() error: %v", err)
	}
	if dir == "" {
		t.Fatal("templateSetDir() returned empty string")
	}
	if !filepath.IsAbs(dir) {
		t.Errorf("templateSetDir() = %q; want absolute path", dir)
	}
	homeDir, _ := os.UserHomeDir()
	expected := filepath.Join(homeDir, ".gouno", "templates")
	if dir != expected {
		t.Errorf("templateSetDir() = %q; want %q", dir, expected)
	}
}

func TestTemplateListEmpty(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := templateListCmd.RunE(cmd, nil)
	if err != nil {
		t.Fatalf("templateListCmd.RunE() error: %v", err)
	}
}

func TestTemplateInstallFromLocalPath(t *testing.T) {
	// Create a temporary template source
	srcDir := t.TempDir()
	os.MkdirAll(filepath.Join(srcDir, "templates"), 0755)
	os.WriteFile(filepath.Join(srcDir, "templates", "test.tmpl"), []byte("hello {{.Name}}"), 0644)

	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.Flags().BoolP("force", "f", false, "")

	// Override templateSetDir for test isolation
	homeDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", origHome)

	err := templateInstallCmd.RunE(cmd, []string{"test-set", srcDir})
	if err != nil {
		t.Fatalf("templateInstallCmd.RunE() error: %v", err)
	}

	// Verify installed
	destPath := filepath.Join(homeDir, ".gouno", "templates", "test-set", "test.tmpl")
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		t.Errorf("expected template file at %s", destPath)
	}
}

func TestTemplateInstallDuplicate(t *testing.T) {
	srcDir := t.TempDir()
	os.WriteFile(filepath.Join(srcDir, "test.tmpl"), []byte("content"), 0644)

	homeDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", origHome)

	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.Flags().BoolP("force", "f", false, "")

	// First install
	err := templateInstallCmd.RunE(cmd, []string{"dup-set", srcDir})
	if err != nil {
		t.Fatalf("first install error: %v", err)
	}

	// Second install without --force should fail
	err = templateInstallCmd.RunE(cmd, []string{"dup-set", srcDir})
	if err == nil {
		t.Error("expected error on duplicate install without --force")
	}
}

func TestTemplateRemoveNonExistent(t *testing.T) {
	homeDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", origHome)

	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := templateRemoveCmd.RunE(cmd, []string{"nonexistent"})
	if err == nil {
		t.Error("expected error when removing non-existent template set")
	}
}

func TestTemplateRemoveExisting(t *testing.T) {
	homeDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", homeDir)
	defer os.Setenv("HOME", origHome)

	// Create a template set directory
	setDir := filepath.Join(homeDir, ".gouno", "templates", "to-remove")
	os.MkdirAll(setDir, 0755)

	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)

	err := templateRemoveCmd.RunE(cmd, []string{"to-remove"})
	if err != nil {
		t.Fatalf("templateRemoveCmd.RunE() error: %v", err)
	}

	if _, err := os.Stat(setDir); !os.IsNotExist(err) {
		t.Error("template set directory should have been removed")
	}
}

func TestCopyDir(t *testing.T) {
	srcDir := t.TempDir()
	destDir := t.TempDir()

	os.WriteFile(filepath.Join(srcDir, "a.txt"), []byte("hello"), 0644)
	os.MkdirAll(filepath.Join(srcDir, "sub"), 0755)
	os.WriteFile(filepath.Join(srcDir, "sub", "b.txt"), []byte("world"), 0644)

	err := copyDir(srcDir, destDir)
	if err != nil {
		t.Fatalf("copyDir() error: %v", err)
	}

	a, err := os.ReadFile(filepath.Join(destDir, "a.txt"))
	if err != nil {
		t.Fatalf("read a.txt: %v", err)
	}
	if string(a) != "hello" {
		t.Errorf("a.txt = %q; want %q", string(a), "hello")
	}

	b, err := os.ReadFile(filepath.Join(destDir, "sub", "b.txt"))
	if err != nil {
		t.Fatalf("read sub/b.txt: %v", err)
	}
	if string(b) != "world" {
		t.Errorf("sub/b.txt = %q; want %q", string(b), "world")
	}
}
