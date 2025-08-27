package gouno

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
)

type TemplateData struct {
	ModulePath  string
	ProjectName string
	RepoURL     string
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
			modulePath = "github.com/" + os.Getenv("USER") + "/" + projectName
		}

		data := TemplateData{
			ModulePath:  modulePath,
			ProjectName: projectName,
			RepoURL:     fmt.Sprintf("https://github.com/%s/%s.git", os.Getenv("USER"), projectName),
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
	newCmd.Flags().StringP("template", "t", "./templates", "Path to the template directory")
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
