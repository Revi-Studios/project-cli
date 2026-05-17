package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
)

var untag bool

func init() {
	rootCmd.AddCommand(tag)
	tag.Flags().BoolVarP(&untag, "untag", "u", false, "Untag the specified tag from the project")
}

var tag = &cobra.Command{
	Use:   "tag <project> <tags>",
	Short: "tag project",
	Long:  "Manage tags on a project",
	Args:  cobra.MinimumNArgs(1),
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

		if len(args) == 1 {
			tags, err := lib.GetTags(path)

			if err != nil {
				return fmt.Errorf("reading existing tags: %w", err)
			}

			fmt.Println(strings.Join(tags, ", "))
			return nil
		}

		argTags := make([]string, 0, len(args)-1)
		for _, tag := range args[1:] {
			argTags = append(argTags, lib.Config.Shorthands.Tag(tag))
		}

		tags, err := lib.GetTags(path)

		if err != nil {
			return fmt.Errorf("reading existing tags: %w", err)
		}

		switch true {
		// If it should untag
		case untag && len(argTags) == 1:
			filtered := make([]string, 0, len(tags))
			for _, tag := range tags {
				if tag != argTags[0] {
					filtered = append(filtered, tag)
				}
			}

			if err := lib.SetTags(path, filtered...); err != nil {
				return fmt.Errorf("adding tag: %w", err)
			}

			fmt.Println("Removed tag:", argTags[0])

		case untag:
			rmTagsMap := make(map[string]struct{}, len(tags))
			for _, tag := range argTags {
				rmTagsMap[tag] = struct{}{}
			}
			filtered := make([]string, 0, len(tags))
			for _, tag := range tags {
				if _, found := rmTagsMap[tag]; !found {
					filtered = append(filtered, tag)
				}
			}

			if err := lib.SetTags(path, filtered...); err != nil {
				return fmt.Errorf("adding tag: %w", err)
			}

			fmt.Println("Removed tags:", strings.Join(argTags, ", "))

		// If there is only one tag to add
		case len(argTags) == 1:
			if err := lib.SetTags(path, append(argTags, tags...)...); err != nil {
				return fmt.Errorf("adding tag: %w", err)
			}

			fmt.Println("Added tag:", argTags[0])

		default:
			if err := lib.SetTags(path, append(argTags, tags...)...); err != nil {
				return fmt.Errorf("adding tag: %w", err)
			}

			fmt.Println("Added tags:", strings.Join(argTags, ", "))

		}

		return nil
	},
}
