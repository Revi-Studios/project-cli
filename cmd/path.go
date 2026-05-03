package cmd

import (
	"fmt"

	"github.com/Revi-Studios/project/lib"
	"github.com/spf13/cobra"
)

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show the path to the project folder",
	Long:  "Show the path to the project folder",
	Args:  cobra.ExactArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := lib.GetConfig()

		if err != nil {
			return fmt.Errorf("getting the config: %w", err)
		}

		fmt.Println(config.ProjectFolderPath)
		return nil
	},
}

var pathSetCmd = &cobra.Command{
	Use:   "set <path>",
	Short: "Set the path to the project folder",
	Long:  "Set the path to the project folder",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		config, err := lib.GetConfig()

		if err != nil && err.Error() != "config file doesn't have any settings" {
			return fmt.Errorf("getting the config: %w", err)
		}

		config.ProjectFolderPath = args[0]

		if err := lib.SaveConfig(config); err != nil {
			return fmt.Errorf("saving config file: %w", err)
		}

		fmt.Println("Project path set:", config.ProjectFolderPath)

		return nil
	},
}

var pathConfig = &cobra.Command{
	Use:   "config",
	Short: "Show the path to the project folder",
	Long:  "Show the path to the project folder",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(lib.ConfigPath)
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)
	pathCmd.AddCommand(pathSetCmd)
	pathCmd.AddCommand(pathConfig)
}
