package gouno

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/spf13/cobra"
)

type TemplateData struct {
	ModulePath  string
	ProjectName string
}

var projectNameRegex = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*$`)

const defaultTemplateRepo = "https://github.com/rushairer/gouno-template"

var runExternalCommand = func(dir, name string, args ...string) error {
	externalCmd := exec.Command(name, args...)
	externalCmd.Dir = dir
	externalCmd.Stdout = os.Stdout
	externalCmd.Stderr = os.Stderr
	return externalCmd.Run()
}

// cloneTemplate clones a template repository into dest. An empty ref follows
// the repository's default branch; otherwise the given immutable branch or tag
// is checked out with a shallow clone.
func cloneTemplate(repo, ref, dest string) error {
	cloneArgs := []string{"clone"}
	if ref != "" {
		cloneArgs = append(cloneArgs, "--branch", ref)
	}
	cloneArgs = append(cloneArgs, "--depth", "1", repo, dest)
	return runExternalCommand("", "git", cloneArgs...)
}

// validateProjectName 校验项目名：合法的目录名，不含路径穿越字符
func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	if strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("project name cannot contain path separators: %s", name)
	}
	if strings.Contains(name, "..") {
		return fmt.Errorf("project name cannot contain '..': %s", name)
	}
	if !projectNameRegex.MatchString(name) {
		return fmt.Errorf("project name must be a valid identifier (letters, digits, underscores, hyphens): %s", name)
	}
	return nil
}

// validateModulePath checks that a Go module path contains only valid characters.
// This prevents YAML/template injection when the module path is rendered into
// config files and source code.
func validateModulePath(path string) error {
	if path == "" {
		return fmt.Errorf("module path cannot be empty")
	}
	if strings.ContainsAny(path, "\n\r\t") {
		return fmt.Errorf("module path must not contain newlines or tabs")
	}
	// Go module paths: letters, digits, dots, slashes, hyphens, underscores
	for _, c := range path {
		if !unicode.IsLetter(c) && !unicode.IsDigit(c) && c != '/' && c != '.' && c != '-' && c != '_' {
			return fmt.Errorf("module path contains invalid character %q", c)
		}
	}
	// 每个路径段不能为空、"." 或 "..",否则 go mod tidy 会失败
	for _, elem := range strings.Split(path, "/") {
		if elem == "" || elem == "." || elem == ".." {
			return fmt.Errorf("module path contains invalid path element %q", elem)
		}
	}
	return nil
}

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new web project from gouno-template",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		projectName := args[0]
		if err := validateProjectName(projectName); err != nil {
			return err
		}

		modulePath, _ := cmd.Flags().GetString("module")
		templateDir, _ := cmd.Flags().GetString("template")
		templateRef, _ := cmd.Flags().GetString("template-ref")
		skipTidy, _ := cmd.Flags().GetBool("skip-tidy")

		if modulePath == "" {
			modulePath = projectName
		}
		if err := validateModulePath(modulePath); err != nil {
			return err
		}

		// Handle template directory logic
		if strings.HasPrefix(templateDir, "git@") || strings.HasPrefix(templateDir, "https://") {
			tempDir, err := os.MkdirTemp("", "gouno-template-")
			if err != nil {
				return fmt.Errorf("creating temporary directory: %w", err)
			}
			defer func() { _ = os.RemoveAll(tempDir) }()

			fmt.Printf("Cloning template from %s to %s\n", templateDir, tempDir)
			if err := cloneTemplate(templateDir, templateRef, tempDir); err != nil {
				return fmt.Errorf("cloning template repository: %w", err)
			}
			templateDir = tempDir
		} else if templateDir == "./templates" {
			if _, err := os.Stat("./templates"); os.IsNotExist(err) {
				tempDir, err := os.MkdirTemp("", "gouno-template-")
				if err != nil {
					return fmt.Errorf("creating temporary directory: %w", err)
				}
				defer func() { _ = os.RemoveAll(tempDir) }()

				fmt.Printf("Local templates directory not found, cloning default template %s to %s\n", defaultTemplateRepo, tempDir)
				if err := cloneTemplate(defaultTemplateRepo, templateRef, tempDir); err != nil {
					return fmt.Errorf("cloning template repository: %w", err)
				}
				templateDir = tempDir
			} else {
				fmt.Printf("Using local templates directory: ./templates\n")
			}
		} else {
			if _, err := os.Stat(templateDir); os.IsNotExist(err) {
				return fmt.Errorf("template directory '%s' does not exist", templateDir)
			}
			fmt.Printf("Using local template directory: %s\n", templateDir)
		}

		data := TemplateData{
			ModulePath:  modulePath,
			ProjectName: projectName,
		}

		fmt.Printf("Creating new project '%s' with module path '%s' from template '%s'\n", projectName, modulePath, templateDir)

		destDir := filepath.Join(".", projectName)
		if _, err := os.Stat(destDir); err == nil {
			return fmt.Errorf("directory %q already exists, refusing to overwrite", destDir)
		}
		if err := copyTemplate(templateDir, destDir, data); err != nil {
			_ = os.RemoveAll(destDir) // 清理已创建的部分文件
			return fmt.Errorf("creating project: %w", err)
		}

		if !skipTidy {
			fmt.Printf("Tidying Go modules...\n")
			if err := tidyProject(destDir); err != nil {
				_ = os.RemoveAll(destDir)
				return err
			}
		}

		fmt.Printf("Project '%s' created successfully in '%s'\n", projectName, destDir)
		fmt.Printf("Next steps:\n")
		fmt.Printf("  1. cd %s\n", projectName)
		fmt.Printf("  2. make dev\n")
		fmt.Printf("  3. Open http://localhost:8080 in your browser\n")
		fmt.Printf("  4. Start coding!\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("module", "m", "", "Go module path (e.g., github.com/your/project)")
	newCmd.Flags().StringP("template", "t", "./templates", fmt.Sprintf("Path to the template directory (default will clone from %s)", defaultTemplateRepo))
	newCmd.Flags().String("template-ref", "", "Immutable branch or tag for a remote template (default: follows the template's default branch)")
	newCmd.Flags().Bool("skip-tidy", false, "Skip running go mod tidy after project creation")
}

// shouldSkipFile 判断是否跳过该文件/目录（检查路径中所有组件）
func shouldSkipFile(relPath string) bool {
	skipNames := map[string]bool{
		".git":      true,
		".idea":     true,
		".DS_Store": true,
		"bin":       true,
		"templates": true,
		".env":      true,
	}
	parts := strings.Split(relPath, string(filepath.Separator))
	for _, part := range parts {
		if skipNames[part] || strings.HasSuffix(part, ".local.yaml") || strings.HasPrefix(part, ".env.") {
			return true
		}
	}
	return false
}

// isRenderableFile 判断文件是否需要进行模板渲染（仅含 {{ 的文本文件）
func isRenderableFile(content string) bool {
	return strings.Contains(content, "{{")
}

func copyTemplate(src, dest string, data TemplateData) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		destPath := filepath.Join(dest, relPath)

		if shouldSkipFile(relPath) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if info.IsDir() {
			return os.MkdirAll(destPath, 0o755)
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		contentBytes, err := io.ReadAll(srcFile)
		closeErr := srcFile.Close()
		if err != nil {
			return err
		}
		if closeErr != nil {
			return closeErr
		}
		content := string(contentBytes)

		// 仅对包含 {{ 的文件进行模板渲染，其余直接复制
		var output string
		if isRenderableFile(content) {
			tmpl, err := template.New("file").Parse(content)
			if err != nil {
				// 无法解析为模板（如 Go 代码中的 {{），直接复制原文
				output = content
			} else {
				var buf strings.Builder
				if err := tmpl.Execute(&buf, data); err != nil {
					return fmt.Errorf("rendering template in %s: %w", relPath, err)
				}
				output = buf.String()
			}
		} else {
			output = content
		}

		// 保留源文件可执行位（如模板中的 run.sh），但去掉 group/other 写位，
		// 避免复制出 0777 这类过度宽松的权限。
		if err := os.WriteFile(destPath, []byte(output), info.Mode().Perm()&0o755); err != nil {
			return err
		}
		return nil
	})
}

func tidyProject(destDir string) error {
	if err := runExternalCommand(destDir, "go", "mod", "tidy"); err != nil {
		return fmt.Errorf("running go mod tidy: %w", err)
	}
	return nil
}
