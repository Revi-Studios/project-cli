package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	Root.CompletionOptions.HiddenDefaultCmd = true
}

var Root = &cobra.Command{
	Use:     "project",
	Short:   "Base command",
	Long:    `A command-line utility for Project managment.`,
	Example: "project\n\nproject list\n\nproject quick",
}

func Execute() {
	if err := Root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
