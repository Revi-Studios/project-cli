package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(tag)
}

var tag = &cobra.Command{
	Use:   "tag <project> <tags>",
	Short: "tag project",
	Long:  "Add a tag to a project",
	Args:  cobra.MinimumNArgs(2),
	ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) > 1 {
			tags := make([]string, 0, len(lib.Config.Tags))
			for _, tag := range lib.Config.Tags {
				tags = append(tags, "\""+tag+"\"")
			}
			return tags, cobra.ShellCompDirectiveNoFileComp
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
		path := filepath.Join(lib.Config.ProjectPath(), name)

		tags := make([]string, 0, len(args)-1)
		for _, tag := range args[1:] {
			tags = append(tags, lib.Config.Shorthands.Tag(tag))
		}
		if err := lib.SetTag(path, tags...); err != nil {
			return fmt.Errorf("adding tag: %w", err)
		}

		return nil
	},
}
