package cmd

import (
	"fmt"
	"os/exec"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
)

var OpenCmd = &cobra.Command{
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
		path := config.ProjectFolderPath + "/" + name + "/"

		// err := exec.Command("osascript", "-e", fmt.Sprintf(`tell application "Terminal" to do script "cd '%s'; clear"`, path)).Run()

		err = exec.Command("open", path).Run()

		if err != nil {
			return fmt.Errorf("opening project %s: %w", name, err)
		}
		fmt.Printf("Opened project: %s\n", name)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(OpenCmd)
}
