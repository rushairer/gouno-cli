package gouno

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// Version is set at build time via -ldflags.
var Version = "dev"

var rootCmd = &cobra.Command{
	Use:   "gouno-cli",
	Short: "gouno-cli is a tool to scaffold Go web projects",
	Long:  "gouno-cli is a tool to scaffold Go web projects from gouno-template.",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of gouno-cli",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("gouno-cli %s\n", Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Version = Version
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
