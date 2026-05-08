package cmd

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
)

var openCmd = &cobra.Command{
	Use:   "open <project-name>",
	Args:  cobra.ExactArgs(1),
	Short: "Open a project",
	Long:  "Open a project",
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := lib.GetConfig()
		if err != nil {
			return fmt.Errorf("getting the config: %w", err)
		}

		name := args[0]
		path := lib.Config.ProjectPath() + "/" + name + "/"

		target, _ := cmd.Flags().GetString("app")
		if target == "" {
			target = config.Defaults.Open
		}

		switch target {
		case "finder", "f":
			err = fmt.Errorf("finder: %w", exec.Command("open", path).Run())
		case "terminal", "t":
			err = fmt.Errorf("terminal: %w", exec.Command("open", "-a", "Terminal", path).Run())
		case "editor", "e":
			err = fmt.Errorf("%s: %w", config.Editor, exec.Command(config.Editor, path).Run())
			target = config.Editor
		default:
			return fmt.Errorf("unknown target: %s", target)
		}

		if errors.Unwrap(err) != nil {
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
	rootCmd.AddCommand(openCmd)
}
