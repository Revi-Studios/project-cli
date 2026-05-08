package cmd

import (
	"fmt"
	"os/exec"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(remove)
}

var remove = &cobra.Command{
	Use:   "remove <name>",
	Short: "Remove a project",
	Long:  "Remove a project",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := lib.Config.ProjectPath() + "/" + name

		if err := exec.Command("trash", path).Run(); err != nil {
			return fmt.Errorf("running 'trash %v': %w", path, err)
		}
		fmt.Println("Project removed at:", path)

		return nil
	},
}
