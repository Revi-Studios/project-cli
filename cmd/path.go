package cmd

import (
	"fmt"
	"project_cli/project/config"

	"github.com/spf13/cobra"
)

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show the path to the project folder",
	Long:  "Show the path to the project folder",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(config.GetConfig().ProjectFolderPath)
	},
}

var pathSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the path to the project folder",
	Long:  "Set the path to the project folder",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: path set <path>")
			return
		}
		localConfig := config.GetConfig()
		localConfig.ProjectFolderPath = args[0]
		if err := config.SaveConfig(localConfig); err != nil {
			fmt.Println("Error writing config file:", err)
			return
		}
		fmt.Println("Project folder path set to:", localConfig.ProjectFolderPath)
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)
	pathCmd.AddCommand(pathSetCmd)
}
