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
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		tags := make([]string, 0, len(lib.Config.Tags))
		for _, tag := range lib.Config.Tags {
			tags = append(tags, "\""+tag+"\"")
		}
		return tags, cobra.ShellCompDirectiveNoFileComp
	},
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		projectPath := lib.Config.ProjectPath() + "/" + name
		err := os.Mkdir(projectPath, 0755)

		if err != nil {
			return fmt.Errorf("creating the directory: %w", err)
		}

		fmt.Println("Project created at:", projectPath)

		if len(args) > 1 {
			tags := make([]string, 0, len(args)-1)
			for _, tag := range args[1:] {
				tags = append(tags, lib.Config.Shorthands.Tag(tag))
			}
			if err := lib.SetTag(projectPath, tags...); err != nil {
				return fmt.Errorf("Error adding tag: %w", err)
			}
			return nil
		}

		if err := lib.SetTag(projectPath, lib.Config.Defaults.Tag); err != nil {
			return fmt.Errorf("Error adding tag: %w", err)
		}

		return nil
	},
}
