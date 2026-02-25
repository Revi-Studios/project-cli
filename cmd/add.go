package cmd

import (
	"fmt"
	"os"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(add)
}

var add = &cobra.Command{
	Use:   "add <name> <tag>",
	Short: "Add a new project",
	Long:  `Add a new project to the project list`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		config, err := lib.GetConfig()

		if err != nil {
			return fmt.Errorf("getting the config: %w", err)
		}

		projectPath := config.ProjectFolderPath + "/" + name
		err = os.Mkdir(projectPath, 0755)

		if err != nil {
			return fmt.Errorf("creating the directory: %w", err)
		}

		fmt.Println("Project created at:", projectPath)

		if err := lib.SetTag(projectPath, args[1]); err != nil {
			return fmt.Errorf("Error adding tag: %w", err)
		}

		return nil
	},
}
