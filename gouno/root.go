package gouno

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gouno-cli",
	Short: "gouno-cli is a tool to generate go code",
	Long: `gouno-cli is a tool to generate go code.
It can generate go code from proto file.`,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing root command: %v", err)
		os.Exit(1)
	}
}
