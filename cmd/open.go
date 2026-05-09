package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:  "open <project-name>",
	Args: cobra.ExactArgs(1),
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
	Short: "Open a project",
	Long:  "Open a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		path := lib.Config.ProjectPath() + "/" + name + "/"

		target, _ := cmd.Flags().GetString("app")
		if target == "" {
			target = lib.Config.Defaults.Open
		}

		var err error

		switch target {
		case "finder", "f":
			err = fmt.Errorf("finder: %w", exec.Command("open", path).Run())
			target = "Finder"
		case "terminal", "t":
			err = fmt.Errorf("terminal: %w", exec.Command("open", "-a", "Terminal", path).Run())
			target = "Terminal"
		case "editor", "e":
			if lib.Config.Editor == "" {
				err = errors.New("editor: not set in config")
				break
			}
			err = fmt.Errorf("%s: %w", lib.Config.Editor, exec.Command(lib.Config.Editor, path).Run())
			target = lib.Config.Editor
		default:
			return fmt.Errorf("unknown target: %s", target)
		}

		if err != nil {
			return fmt.Errorf("opening project %s: %w", name, err)
		}

		// Turns the first character to uppercase
		targetToCapital := func() string { str := []rune(target); return strings.ToUpper(string(str[0])) + string(str[1:]) }()

		fmt.Printf("Opened project: \"%s\" in %s\n", name, targetToCapital)

		return nil
	},
}

func init() {
	openCmd.Flags().StringP("app", "a", "", "--app <platform/target>")
	openCmd.RegisterFlagCompletionFunc("app", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return []string{"finder", "terminal", "editor"}, cobra.ShellCompDirectiveNoFileComp
	})
	rootCmd.AddCommand(openCmd)
}
