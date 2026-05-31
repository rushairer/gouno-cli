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
	if !unicode.IsLetter(rune(name[0])) && name[0] != '_' {
		return fmt.Errorf("project name must start with a letter or underscore: %s", name)
	}
	if !projectNameRegex.MatchString(name) {
		return fmt.Errorf("project name must be a valid identifier (letters, digits, underscores, hyphens): %s", name)
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

		if modulePath == "" {
			modulePath = projectName
		}

		// Handle template directory logic
		if strings.HasPrefix(templateDir, "git@") || strings.HasPrefix(templateDir, "https://") {
			tempDir, err := os.MkdirTemp("", "gouno-template-")
			if err != nil {
				return fmt.Errorf("creating temporary directory: %w", err)
			}
			defer os.RemoveAll(tempDir)

			fmt.Printf("Cloning template from %s to %s\n", templateDir, tempDir)
			gitCmd := exec.Command("git", "clone", templateDir, tempDir)
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
			if err := gitCmd.Run(); err != nil {
				return fmt.Errorf("cloning template repository: %w", err)
			}
			templateDir = tempDir
		} else if templateDir == "./templates" {
			if _, err := os.Stat("./templates"); os.IsNotExist(err) {
				tempDir, err := os.MkdirTemp("", "gouno-template-")
				if err != nil {
					return fmt.Errorf("creating temporary directory: %w", err)
				}
				defer os.RemoveAll(tempDir)

				fmt.Printf("Local templates directory not found, cloning default template from https://github.com/rushairer/gouno-template to %s\n", tempDir)
				gitCmd := exec.Command("git", "clone", "https://github.com/rushairer/gouno-template", tempDir)
				gitCmd.Stdout = os.Stdout
				gitCmd.Stderr = os.Stderr
				if err := gitCmd.Run(); err != nil {
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
		if err := copyTemplate(templateDir, destDir, data); err != nil {
			os.RemoveAll(destDir) // 清理已创建的部分文件
			return fmt.Errorf("creating project: %w", err)
		}

		// 写入 .gouno.yaml 配置
		templateSet, _ := cmd.Flags().GetString("template-set")
		if templateSet != "" {
			cfgContent := fmt.Sprintf("template-set: %s\n", templateSet)
			cfgPath := filepath.Join(destDir, ".gouno.yaml")
			if err := os.WriteFile(cfgPath, []byte(cfgContent), 0644); err != nil {
				return fmt.Errorf("writing .gouno.yaml: %w", err)
			}
			fmt.Printf("Template set '%s' saved to .gouno.yaml\n", templateSet)
		}

		fmt.Printf("Project '%s' created successfully in '%s'\n", projectName, destDir)
		fmt.Printf("Next steps:\n")
		fmt.Printf("  1. cd %s\n", projectName)
		fmt.Printf("  2. go mod tidy\n")
		fmt.Printf("  3. make dev\n")
		fmt.Printf("  4. Open http://localhost:8080 in your browser\n")
		fmt.Printf("  5. Start coding!\n")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("module", "m", "", "Go module path (e.g., github.com/your/project)")
	newCmd.Flags().StringP("template", "t", "./templates", "Path to the template directory (default will clone from https://github.com/rushairer/gouno-template)")
	newCmd.Flags().String("template-set", "", "Template set name for code generation (saved to .gouno.yaml)")
}

// shouldSkipFile 判断是否跳过该文件/目录（检查路径中所有组件）
func shouldSkipFile(relPath string) bool {
	skipNames := map[string]bool{
		".git":      true,
		".idea":     true,
		".DS_Store": true,
		"bin":       true,
		"templates": true,
	}
	parts := strings.Split(relPath, string(filepath.Separator))
	for _, part := range parts {
		if skipNames[part] {
			return true
		}
		if strings.HasPrefix(part, ".git") {
			return true
		}
		if strings.HasSuffix(part, ".tmpl") {
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
			return os.MkdirAll(destPath, info.Mode())
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		contentBytes, err := io.ReadAll(srcFile)
		if err != nil {
			return err
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

		if err := os.WriteFile(destPath, []byte(output), info.Mode()); err != nil {
			return err
		}
		return nil
	})
}
