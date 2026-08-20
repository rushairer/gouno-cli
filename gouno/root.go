package gouno

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/spf13/cobra"
)

// Version 可在构建时通过 -ldflags 注入,例如:
//
//	go build -ldflags "-X github.com/rushairer/gouno-cli/gouno.Version=v1.1.0" .
//
// 若未注入,getVersion() 会回退读取 go module 构建信息中的版本。
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
		cmd.Printf("gouno-cli %s\n", getVersion())
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.Version = getVersion()
}

// getVersion 返回版本号,优先级:
//  1. 构建时通过 -ldflags 注入的 Version(如 Makefile / goreleaser 发布流程)
//  2. go module 构建信息中的 main module 版本(go install pkg@version 场景)
//  3. 兜底返回 "dev"
func getVersion() string {
	if Version != "" && Version != "dev" {
		return Version
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if v := info.Main.Version; v != "" && v != "(devel)" {
			return v
		}
	}
	return "dev"
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
