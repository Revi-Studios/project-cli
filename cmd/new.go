package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
)

var open bool

func init() {
	new.Flags().BoolVarP(&open, "open", "o", false, "Open project in default target application")
	Root.AddCommand(new)
}

var new = &cobra.Command{
	Use:     "new <name> <tag>",
	Short:   "Create a new project",
	Long:    `Create a new project in the projects folder`,
	Example: "project new snake-game web,javascript\n\nproject new api-server go,backend\n\nproject new notes",
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		// tags := make([]string, 0, len(lib.Config.Tags))
		// for _, tag := range lib.Config.Tags {
		// 	tags = append(tags, "\""+tag+"\"")
		// }
		tags := make([]string, 0, len(lib.Config.Shorthands.Tags))
		for tag := range lib.Config.Shorthands.Tags {
			tags = append(tags, tag)
		}

		return tags, cobra.ShellCompDirectiveNoFileComp
	},
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		projectPath := lib.Config.ProjectPath() + "/" + name
		err := os.Mkdir(projectPath, 0755)

		if errors.Is(err, os.ErrExist) {
			fmt.Printf("Project \"%s\" already exist\n", name)
			return nil
		}

		if err != nil {
			return fmt.Errorf("creating the directory: %w", err)
		}

		fmt.Printf("Project \"%s\" created at: \033]8;;%s\033\\%s\033]8;;\033\\\n", name, fmt.Sprintf("file://%s/", projectPath), projectPath)

		if len(args) > 1 {
			tags := make([]string, 0, len(args)-1)
			for _, tag := range args[1:] {
				tags = append(tags, lib.Config.Shorthands.Tag(tag))
			}
			if err := lib.SetTags(projectPath, tags...); err != nil {
				return fmt.Errorf("Error adding tag: %w", err)
			}
			return nil
		}

		if err := lib.SetTags(projectPath, lib.Config.Defaults.Tag); err != nil {
			return fmt.Errorf("Error adding tag: %w", err)
		}

		if open {
			lib.OpenProjectInApp(name, "")
		}

		return nil
	},
}
