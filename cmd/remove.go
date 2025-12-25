package cmd

import (
	"fmt"
	"os/exec"
	"project_cli/project/config"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(remove)
}

var remove = &cobra.Command{
	Use:   "remove",
	Short: "Remove a project",
	Long:  "Remove a project",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a project name")
			return
		}
		projectName := args[0]
		projectPath := config.GetConfig().ProjectFolderPath + "/" + projectName
		err := exec.Command("trash", projectPath).Run()
		if err != nil {
			fmt.Println("Error removing project:", err)
			return
		}
		fmt.Println("Project removed:", projectPath)
	},
}
