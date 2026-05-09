package cmd

import (
	"fmt"
	"os"
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
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		folders, err := func() ([]string, error) {
			files, err := os.ReadDir(lib.Config.ProjectPath() + "/")
			if err != nil {
				return nil, err
			}
			folders := make([]string, 0, len(files))
			for _, file := range files {
				if file.IsDir() {
					folders = append(folders, file.Name())
				}
			}
			return folders, nil
		}()

		if err != nil {
			return nil, cobra.ShellCompDirectiveError
		}

		return folders, cobra.ShellCompDirectiveNoFileComp
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		name := lib.Config.Shorthands.Project(args[0])
		path := lib.Config.ProjectPath() + "/" + name

		if err := exec.Command("trash", path).Run(); err != nil {
			return fmt.Errorf("running 'trash %v': %w", path, err)
		}
		fmt.Println("Project removed at:", path)

		return nil
	},
}
