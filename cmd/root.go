package cmd

import (
	"fmt"
	"os"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
)

var Root = &cobra.Command{
	Use:     "project",
	Short:   "Base command",
	Long:    `Base command. If another command isn't added, it opens the project folder.`,
	Example: "project\n\nproject list\n\nproject quick",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := lib.OpenFolder(lib.Config.ProjectPath())
		if err != nil {
			return fmt.Errorf("Faild to open the project folder: %w", err)
		}
		return nil
	},
}

func Execute() {
	if err := Root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
