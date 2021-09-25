package main

import (
	"os"

	"github.com/keboola/keboola-as-code/internal/pkg/cli"
	"github.com/keboola/keboola-as-code/internal/pkg/filesystem/aferofs"
	"github.com/keboola/keboola-as-code/internal/pkg/interaction"
)

func main() {
	// Run command
	prompt := interaction.NewPrompt(os.Stdin, os.Stdout, os.Stderr)
	cmd := cli.NewRootCommand(os.Stdin, os.Stdout, os.Stderr, prompt, aferofs.NewLocalFsFindProjectDir)
	os.Exit(cmd.Execute())
}
