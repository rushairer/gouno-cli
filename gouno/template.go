package gouno

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

const templateDirName = ".gouno"
const templatesDirName = "templates"

func templateSetDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, templateDirName, templatesDirName), nil
}

// validateTemplateName checks that a template name does not contain path
// separators or traversal sequences that could escape the templates directory.
func validateTemplateName(name string) error {
	if name == "" {
		return fmt.Errorf("template name must not be empty")
	}
	if strings.ContainsAny(name, "/\\") {
		return fmt.Errorf("template name must not contain path separators")
	}
	if name == "." || name == ".." || strings.HasPrefix(name, "..") {
		return fmt.Errorf("template name must not be a path traversal sequence")
	}
	return nil
}

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage template sets",
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed template sets",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir, err := templateSetDir()
		if err != nil {
			return err
		}
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			cmd.Println("No template sets installed.")
			cmd.Printf("Install with: gouno-cli template install <name> <git-url>\n")
			return nil
		}

		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to read template directory: %w", err)
		}

		if len(entries) == 0 {
			cmd.Println("No template sets installed.")
			return nil
		}

		cmd.Printf("Installed template sets (%s):\n", dir)
		for _, entry := range entries {
			if entry.IsDir() {
				cmd.Printf("  - %s\n", entry.Name())
			}
		}
		return nil
	},
}

var templateInstallCmd = &cobra.Command{
	Use:   "install <name> <git-url-or-path>",
	Short: "Install a template set from git URL or local path",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		source := args[1]

		if err := validateTemplateName(name); err != nil {
			return err
		}

		dir, err := templateSetDir()
		if err != nil {
			return err
		}
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create template directory: %w", err)
		}

		destPath := filepath.Join(dir, name)
		if _, err := os.Stat(destPath); err == nil {
			force, _ := cmd.Flags().GetBool("force")
			if !force {
				return fmt.Errorf("template set %q already exists (use --force to overwrite)", name)
			}
			os.RemoveAll(destPath)
		}

		// 判断来源：本地路径 or git URL
		if info, err := os.Stat(source); err == nil && info.IsDir() {
			// 本地目录：优先使用 templates/ 子目录
			tmplSubdir := filepath.Join(source, "templates")
			if _, err := os.Stat(tmplSubdir); err == nil {
				source = tmplSubdir
			}
			if err := copyDir(source, destPath); err != nil {
				return fmt.Errorf("failed to copy template: %w", err)
			}
			cmd.Printf("Template set %q installed from %s\n", name, source)
		} else {
			// git clone
			tempDir, err := os.MkdirTemp("", "gouno-install-")
			if err != nil {
				return fmt.Errorf("failed to create temp directory: %w", err)
			}
			defer os.RemoveAll(tempDir)

			gitCmd := exec.Command("git", "clone", source, tempDir)
			gitCmd.Stdout = os.Stdout
			gitCmd.Stderr = os.Stderr
			if err := gitCmd.Run(); err != nil {
				return fmt.Errorf("failed to clone template: %w", err)
			}
			// 优先使用 templates/ 子目录
			tmplSubdir := filepath.Join(tempDir, "templates")
			if _, err := os.Stat(tmplSubdir); err == nil {
				tempDir = tmplSubdir
			}
			if err := copyDir(tempDir, destPath); err != nil {
				return fmt.Errorf("failed to copy template: %w", err)
			}
			cmd.Printf("Template set %q installed from %s\n", name, source)
		}

		cmd.Printf("Location: %s\n", destPath)
		return nil
	},
}

var templateRemoveCmd = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove an installed template set",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		if err := validateTemplateName(name); err != nil {
			return err
		}
		dir, err := templateSetDir()
		if err != nil {
			return err
		}

		destPath := filepath.Join(dir, name)
		if _, err := os.Stat(destPath); os.IsNotExist(err) {
			return fmt.Errorf("template set %q not found", name)
		}

		if err := os.RemoveAll(destPath); err != nil {
			return fmt.Errorf("failed to remove template set: %w", err)
		}
		cmd.Printf("Template set %q removed\n", name)
		return nil
	},
}

func init() {
	templateInstallCmd.Flags().BoolP("force", "f", false, "overwrite existing template set")
	templateCmd.AddCommand(templateListCmd, templateInstallCmd, templateRemoveCmd)
	rootCmd.AddCommand(templateCmd)
}

// copyDir 递归复制目录
func copyDir(src, dest string) error {
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
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return os.WriteFile(destPath, data, info.Mode())
	})
}
