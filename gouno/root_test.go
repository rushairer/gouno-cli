package gouno

import (
	"bytes"
	"strings"
	"testing"

	"github.com/spf13/cobra"
)

func TestVersionCmd(t *testing.T) {
	buf := new(bytes.Buffer)
	cmd := &cobra.Command{}
	cmd.SetOut(buf)

	versionCmd.Run(cmd, nil)

	out := buf.String()
	if !strings.Contains(out, "gouno-cli ") {
		t.Errorf("version output = %q; want prefix 'gouno-cli '", out)
	}
	if !strings.Contains(out, getVersion()) {
		t.Errorf("version output = %q; want contains version %q", out, getVersion())
	}
}

func TestRootCmdSilenceErrors(t *testing.T) {
	// 关闭 cobra 默认错误打印,由 Execute() 统一输出一次,避免双重打印
	if !rootCmd.SilenceErrors {
		t.Error("rootCmd.SilenceErrors should be true to avoid double error printing")
	}
	if !rootCmd.SilenceUsage {
		t.Error("rootCmd.SilenceUsage should be true to avoid usage spam on errors")
	}
}

func TestRootCmdSubcommands(t *testing.T) {
	// 根命令应挂载 new / version 子命令
	names := make(map[string]bool)
	for _, c := range rootCmd.Commands() {
		names[c.Name()] = true
	}
	for _, want := range []string{"new", "version"} {
		if !names[want] {
			t.Errorf("root command missing subcommand %q; have %v", want, names)
		}
	}
}
