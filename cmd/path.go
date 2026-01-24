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
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(lib.GetConfig().ProjectFolderPath)
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
		localConfig := lib.GetConfig()
		localConfig.ProjectFolderPath = args[0]
		if err := lib.SaveConfig(localConfig); err != nil {
			fmt.Println("Error writing config file:", err)
			return
		}
		fmt.Println("Project folder path set to:", localConfig.ProjectFolderPath)
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
