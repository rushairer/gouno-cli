package gouno

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type TemplateData struct {
	ModulePath  string
	ProjectName string
}

var newCmd = &cobra.Command{
	Use:   "new [project-name]",
	Short: "Create a new web project from go-uno template",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		modulePath, _ := cmd.Flags().GetString("module")
		templateDir, _ := cmd.Flags().GetString("template")

		if modulePath == "" {
			modulePath = projectName
		}

		// If templateDir is not provided, clone from default repository
		if templateDir == "./templates" {
			tempDir, err := os.MkdirTemp("", "gouno-template-")
			if err != nil {
				fmt.Printf("Error creating temporary directory: %v\n", err)
				os.Exit(1)
			}
			defer os.RemoveAll(tempDir) // Clean up the temporary directory

			fmt.Printf("Cloning default template from https://github.com/rushairer/gouno-template to %s\n", tempDir)
			cmd := exec.Command("git", "clone", "https://github.com/rushairer/gouno-template", tempDir)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Error cloning template repository: %v\n", err)
				os.Exit(1)
			}
			templateDir = tempDir
		}

		data := TemplateData{
			ModulePath:  modulePath,
			ProjectName: projectName,
		}

		fmt.Printf("Creating new project '%s' with module path '%s' from template '%s'\n", projectName, modulePath, templateDir)

		destDir := filepath.Join(".", projectName)
		err := copyTemplate(templateDir, destDir, data)
		if err != nil {
			fmt.Printf("Error creating project: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Project '%s' created successfully in '%s'\n", projectName, destDir)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().StringP("module", "m", "", "Go module path (e.g., github.com/your/project)")
	newCmd.Flags().StringP("template", "t", "./templates", "Path to the template directory (default will clone from https://github.com/rushairer/gouno-template)")
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

		if info.IsDir() {
			return os.MkdirAll(destPath, info.Mode())
		}

		// Skip certain files/directories
		if strings.Contains(relPath, ".git") ||
			strings.Contains(relPath, ".idea") ||
			strings.Contains(relPath, ".DS_Store") ||
			strings.Contains(relPath, "bin") ||
			strings.Contains(relPath, "templates") {
			return nil
		}

		srcFile, err := os.Open(path)
		if err != nil {
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			return err
		}
		defer destFile.Close()

		// Read content and apply template rendering
		contentBytes, err := io.ReadAll(srcFile)
		if err != nil {
			return err
		}
		content := string(contentBytes)

		// Apply template rendering for placeholders
		tmpl, err := template.New("file").Parse(content)
		if err != nil {
			// If it's not a template, just copy as is
			_, err = destFile.WriteString(content)
			return err
		}

		var buf strings.Builder
		err = tmpl.Execute(&buf, data)
		if err != nil {
			return err
		}

		_, err = destFile.WriteString(buf.String())
		return err
	})
}
