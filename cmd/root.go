package cmd

import (
	"context"
	"log/slog"
	"os"

	"github.com/adriancmiranda-wm/cli/internal/version"
	"github.com/charmbracelet/fang"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "wm",
	Short: "WM - CLI scaffolding tool",
	Long:  "WM - CLI scaffolding tool",
	Example: `
	# Run in interactive mode
	wm

	# Run with debug logging
	wm -d

	# View the version
	wm -v

	# Run in dangerous mode (auto-accept all prompts)
	wm -y

	# Run with debug logging and specify the project directory
	wm -d -c /path/to/project

	# Generate template
	wm template init -t golib -n MeuProj -a "Adrian"
	`,
}

func Execute() {
	if err := fang.Execute(
		context.Background(),
		rootCmd,
		fang.WithVersion(version.Version),
		fang.WithNotifySignal(os.Interrupt),
	); err != nil {
		slog.Error("Failed to execute fang command", "error", err)
		os.Exit(1)
	}
}
