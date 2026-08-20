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
	// 关闭 cobra 的默认错误打印,统一由 Execute() 输出一次
	SilenceErrors: true,
	SilenceUsage:  true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version of gouno-cli",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Printf("gouno-cli %s\n", Version)
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
